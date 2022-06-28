#!/usr/bin/env python
"""
Positional parameters:

1. number of df row to begin with
2. number of last df row to clean
3. path to df for cleaning
4. file name to save labels
"""

import os
import time
import numpy as np
import pandas as pd
import sys

from scipy.spatial import distance
from sklearn.metrics import classification_report
import transformers

from transformers import pipeline
import sentence_transformers 
from sentence_transformers import SentenceTransformer

from collections import defaultdict
from bertopic import BERTopic

### Zero Shot Sentence-BERT based filtration

def classifier_semantic(text, targets, model):
    target_emb = model.encode(targets)
    text_emb = model.encode(text)
    
    probs = [ 1 - distance.cosine(text_emb, t) for t in target_emb]
    probs = probs / sum(probs)
    
    response = dict()
    for p, t in zip(probs, targets):
        response[t] = p 
    
    temp = dict(sorted(response.items(), key=lambda item: item[1], reverse=True))
    response = dict()
    response['sequence'] = text
    response['scores'] = list(temp.values())
    response['labels'] = list(temp.keys())
    
    return response

def clean_dataset(df, candidate_labels, classifier, cand_adv_index, debug = 100000, threshold = 0):
    labels = list()
    semantic_classifier_model = SentenceTransformer('sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2')
    for index, row in df.iterrows():
        if index % debug == 0:
            print('[ZS cleaning debug] index = ', index)    
        res = classifier(row['caption'], candidate_labels, semantic_classifier_model)
        #print(res)
        if res['scores'][0] > threshold:
            val = 0
            for ind in cand_adv_index:
                if res['labels'][0] == candidate_labels[ind]: # equals to noise label (food, ...)
                    val = 1
                    break
            labels.append(val)
        else:
            labels.append(1) # not enough confidence to label as event
    return labels

### BERTopic based filtration

VALID_DF = os.environ['STORAGE'] + 'instabert/filtering_data/labeled_data_full_union.csv'
BERT_PREDICTOR = os.environ['STORAGE'] + 'instabert/bert_topic/nyc_all_200_part_1_10'

def make_bertopic_predictor(model_path: str):
    topic_model = BERTopic.load(model_path)
    return lambda texts: topic_model.transform(texts)[0] 

def prepare_valid_df(path: str) -> (pd.DataFrame):
    df = pd.read_csv(path)
    df = df[df.label != "trash"]
    df.label = df.label.apply(lambda l: l.replace('#other', '').replace('_other', 'TEMP').replace('other#', '').replace('TEMP', '_other'))
    df.reset_index(drop=True, inplace=True)
    counts = df.iloc[:, 3:-1].astype(int).sum(axis=0)
    return df, counts

def topic_predict(df: pd.DataFrame(), predictor, text_col: str = "caption", predicted_label: str = "predicted") -> (pd.DataFrame, pd.Series):
    df[predicted_label] = predictor(df[text_col])
    return df

def topic_dominant_labels(label_df: pd.DataFrame, label_counts: pd.Series, predicted_label: str = "predicted") -> dict:
    topic_to_labels_influence = defaultdict(lambda: defaultdict(int))
 
    
    # find counts of each labels in topics  
    for _, row in label_df.iterrows():
        for gold in row.label.split("#"):
            topic_to_labels_influence[row[predicted_label]][gold] += 1
            
    
    topic_class = defaultdict(list) 
    for topic_id in topic_to_labels_influence:
        topic = topic_to_labels_influence[topic_id]
        
        # calc influences of labels on the topic
        sum_labels = 0
        for label in topic:
            topic[label] /= label_counts[label]
            sum_labels += topic[label]
    
        # choose minimum of dominant topic
        topic_founded = False
        for i in range(1, len(topic) + 1):
            for label in topic:
                if topic[label] >= sum_labels / i:
                    topic_class[topic_id].append((label, topic[label]))
                    topic_founded = True
            if topic_founded:
                break
    return topic_class

def cluster_to_bool_val(code, code_dict, adv_label_list ):
    if code == -1: # could not determine topic, then do not filter it
        return 0
    if code not in code_dict:
        return 0
    xs = np.argmax([x[1] for x in topic_class[code]])
    if code_dict[code][xs][0] in adv_label_list:
        return 1
    return 0

if __name__ == '__main__':
    start_time = time.time()
    try:
        start = int(sys.argv[1])
        finish = int(sys.argv[2])
            
        data_path = sys.argv[3] + str(start) + '_' + str(finish) + '.csv'
        save_path = sys.argv[4]
        
        print('Subprocess starts:', start, ' - ', finish, '\tPath to data:', data_path, '\tFile to save:', save_path)
    except:
        print('Subprocess: bad parameters')
        sys.exit(1)
    
    df = pd.read_csv(data_path, lineterminator='\n')
    df.caption = df.caption.apply(lambda x: x.replace('\n', ' ') if type(x) != float else print(x))


    print('Cleaning with zero-shot')
    candidate_labs = [
        'other',
        'food',
        'advertisement',
        'spam',
        'promotion',
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
    
    y_pred = clean_dataset(df, candidate_labs, classifier_semantic,[0, 1, 2, 3, 4])
    print('ZS is finished ')
    np.save(save_path + '/zero_' + str(start) + '_' + str(finish), np.array(y_pred))
    
    print('Cleaning with Bertopic')
    #prepare BERTopic
    bt_predictor = make_bertopic_predictor(BERT_PREDICTOR)
    valid_df, label_counts = prepare_valid_df(VALID_DF)
    valid_df = topic_predict(valid_df, bt_predictor, predicted_label="bertopic_all_200_part_10")
    topic_class = topic_dominant_labels(valid_df, label_counts, "bertopic_all_200_part_10") 
    adver_labs = ['adv_event', 'adv_other', 'food', 'other', 'retrospective_event', 'future_event' ]
    
    #clean
    df = topic_predict(df, bt_predictor, predicted_label="bertopic")
    df['adv_label2'] = df.bertopic.apply(cluster_to_bool_val, args=[topic_class, adver_labs])
    print('Bertopicis finished ')
    np.save(save_path + '/bertopic_' + str(start) + '_' + str(finish), np.array(df['adv_label2'].to_list()))
    
    print(f"Filtration suboridnate script has finished in {(time.time() - start_time)/60} minutes")

