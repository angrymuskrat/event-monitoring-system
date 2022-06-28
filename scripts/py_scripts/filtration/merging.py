#!/usr/bin/env python
"""
This script calculates similarity between two rows which should be used later for merging purposes.

Note: this script is only created to scratch the main idea of similarity calculation. 
It takes too much time to initialize all necessary methods, therefore the script is not supposed to be used directly.

If determined topics of posts are different and are not from the same list (both public or both private), then only hastags and semantic criteria are used.

Positional parameters:

1. path to a csv table
2. index of the first row for similarity calculation
3. index of the second row for similarity calculation
"""

import os
import sys
import re
from datetime import datetime
import numpy as np
import pandas as pd

import transformers
from transformers import pipeline
import sentence_transformers 
from sentence_transformers import SentenceTransformer

from scipy.spatial import distance
from sklearn.metrics.pairwise import haversine_distances
from math import radians

candidate_labs = [
'music concert', 
'exhibition', 
'festival',
'conference',
'calendar holiday',
'sport event',
'flashmob', 
'accident',
'stroll walking',
'wedding birthday',
'private event',
'public event']

candidates_time = dict() # duration in days
# public events
candidates_time['music concert'] = 2
candidates_time['exhibition'] = 90
candidates_time['festival'] = 30
candidates_time['conference'] = 10
candidates_time['calendar holiday'] = 2
candidates_time['sport event'] = 2
candidates_time['flashmob'] = 2
candidates_time['accident'] = 2
candidates_time['public event'] = 7
# private events
candidates_time['stroll walking'] = 2
candidates_time['wedding birthday'] = 2
candidates_time['private event'] = 2

candidates_space = dict() # radius in meters
# public events
candidates_space['music concert'] = 100
candidates_space['exhibition'] = 100
candidates_space['festival'] = 400
candidates_space['conference'] = 100
candidates_space['calendar holiday'] = 10000
candidates_space['sport event'] = 1000
candidates_space['flashmob'] = 100
candidates_space['accident'] = 200
candidates_space['public event'] = 1000
# private events
candidates_space['stroll walking'] = 100
candidates_space['wedding birthday'] = 50
candidates_space['private event'] = 50

private_events = ['stroll walking', 'wedding birthday', 'private event']

semantic_classifier_model = SentenceTransformer('sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2')
target_emb = semantic_classifier_model.encode(candidate_labs)

def connect_semantic(text1, text2):
    text1_emb = semantic_classifier_model.encode(text1)
    text2_emb = semantic_classifier_model.encode(text2)
    return 1 - distance.cosine(text1_emb, text2_emb)

def connect_time(time1, time2, topic1, topic2, debug=False):
    # check wich cluster it is, choose time shift
    if debug == True:
        print(datetime.fromtimestamp(time1))
        print(datetime.fromtimestamp(time2))
        
    date_diff = abs((float(time1)-float(time2))/(60*60*24))
    
    if debug == True:
        print('date diff', date_diff, candidates_time[topic1], topic1, topic2)

    if topic1 == topic2: # if the same topic
        allowed_diff = candidates_time[topic1]
    elif topic1 in private_events and topic2 in private_events:
        if debug == True:
            print('Treat as private events')
        allowed_diff = candidates_time['private event']
    elif topic1 not in private_events and topic2 not in private_events:
        if debug == True:
            print('Treat as public events')
        allowed_diff = candidates_time['public event']
    else:
        return None
        
    if allowed_diff > date_diff:
        return 1
    elif allowed_diff * 2 > date_diff:
        return 0.5
    else:
        return 0

def connect_space(lon1, lat1, lon2, lat2, topic1, topic2, debug=False):
    # check wich cluster it is, choose space shift
    if debug == True:
        print(lon1, lat1)
        print(lon2, lat2)
        
    rad1 = [radians(lon1), radians(lat1)]
    rad2 = [radians(lon2), radians(lat2)]
    
    meters_diff = haversine_distances([rad1, rad2]) * 6371000  # multiply by Earth radius to get km
    meters_diff = meters_diff[0][1]
    
    if debug == True:
        print('meters_diff', meters_diff, candidates_space[topic1], topic1, topic2)
    
    if topic1 == topic2: # if the same topic
        allowed_diff = candidates_space[topic1]
    elif topic1 in private_events and topic2 in private_events:
        if debug == True:
            print('Treat as private events')
        allowed_diff = candidates_space['private event']
    elif topic1 not in private_events and topic2 not in private_events:
        if debug == True:
            print('Treat as public events')
        allowed_diff = candidates_space['public event']
    else:
        return None
    
    if allowed_diff > meters_diff:
        return 1
    elif allowed_diff * 2 > meters_diff:
        return 0.5
    else:
        return 0

def connect_hastags(text1, text2, debug=False):
    hashtag_list1 = set(re.findall("#(\w+)", text1))
    hashtag_list2 = set(re.findall("#(\w+)", text2))
    intersect = hashtag_list1 & hashtag_list2
    if debug:
        print(len(intersect), intersect)
    return np.tanh(len(intersect))

def get_topic(text):
    text_emb = semantic_classifier_model.encode(text)
    probs = [ 1 - distance.cosine(text_emb, t) for t in target_emb]
    probs = probs / sum(probs)
    response = dict()
    for p, t in zip(probs, candidate_labs):
        response[t] = p 
    answer = dict(sorted(response.items(), key=lambda item: item[1], reverse=True))
    topic = list(answer.keys())[0]
    return topic

def get_similarity(row1, row2, debug = False):
    
    weigths = [ 1.1, 1.1, 1.5, 1.0 ]
    
    s1 = connect_time(row1.timestamp, row2.timestamp, row1.topic, row2.topic, debug=debug)
    s2 = connect_space(row1.lat, row1.lon, row2.lat, row2.lon, row1.topic, row2.topic, debug=debug)
    s3 = connect_semantic(row1.caption, row2.caption)
    s4 = connect_hastags(row1.caption, row2.caption, debug=debug)
    
    total = [ weigths[index]*i for index, i in enumerate([s1, s2, s3, s4]) if i != None]
    if debug:
        print('Total=', total)
    return sum(total)/len(total)

if __name__ == '__main__':  
    try:
        data_path = sys.argv[1]
        index1 = int(sys.argv[2])
        index2 = int(sys.argv[3])
    except:
        print('Error: Bad parameters')
        sys.exit(1)

    print('Data path:\t', data_path)
    
    print('------------------------')
    print('Reading df:')
    df = pd.read_csv(data_path)
    df = df.dropna(subset = ['caption']).reset_index(drop=True)
    df = df[:10]
    
    print('------------------------')
    print('Topic obtaining:')
    df['topic'] = df.caption.apply(get_topic)
    
    print('------------------------')
    print('Similarity calculation:')
    get_similarity(df.iloc[index1], df.iloc[index2], debug=True)
    
    sys.exit(0)
