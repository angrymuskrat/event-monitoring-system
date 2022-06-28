#!/usr/bin/env python
"""
Positional parameters:

1. path to a csv table which should be filtered 
2. path to a folder to save resulting csv and temporary files
3. number of processes to clean (take into account gpu memory, one process needs ~1700MiB)
"""

import time
import sys
import os
import math
import pandas as pd
import numpy as np
from subprocess import Popen

try:
    data_path = sys.argv[1]
    save_path = sys.argv[2]
    proc_num = int(sys.argv[3])
except:
    print('Bad parameters')
    syx.exit(1)
    
print('Data path:\t', data_path)
print('Save path:\t', save_path)
print('Proc num:\t', proc_num)

print('------------------------')
print('Reading df:')
df = pd.read_csv(data_path)
df = df.dropna(subset = ['caption']).reset_index(drop=True)

df = df[:63]

total_number = len(df)
shift = math.ceil(total_number/proc_num)
print('Rows number =', total_number,', Shift =', shift)

# Save N ( N = proc_num) new datasets
print('------------------------')
print('Creating new temporary dfs:')
start = finish = index = 0
while( finish != total_number):
    index += 1
    finish = start + shift
    
    if (finish > total_number) or (index == proc_num): # cathes last part, unnecessary double condition
        finish = total_number
    
    temp_df = df[start:finish]
    print('\t', index, ':', start, '-', finish)
    temp_df.to_csv(save_path + 'tempdf_' + str(start) + '_' + str(finish)+'.csv', sep=',')
    start = finish
    
# Call subprcocesses
print('------------------------')
print('Call subprocesses:')
i = 0
start = i * shift
finish = 0
proc_list = list()
while(finish != total_number):
    finish = start + shift
    if finish > total_number:
        finish = total_number
    cmd = [ sys.executable, './filter_subordinate.py', str(start), str(finish), save_path+'tempdf_', save_path,]
    p = Popen(cmd)
    proc_list.append(p) 
    print(f'\t', i+1, ': pid =', p.pid)
    start = finish
    i += 1

# Wait for all processes
for proc in proc_list:
    proc.communicate() 

print('------------------------')
print('Aggregating results of subprocesses:')
# Concatenate results
start = finish = 0
dfs = list()
while( finish != df.shape[0] ):
    finish = start + shift
    if finish > df.shape[0]:
        finish = df.shape[0]
    temp_df = df[start:finish]
    
    np_path = save_path + '/zero_' + str(start) + '_' + str(finish)+'.npy'
    temp_arr = np.load(np_path)
    print('For zero ', str(start), ':',str(finish), temp_df.shape[0],'for', np_path, len(temp_arr))
    assert(temp_df.shape[0] == len(temp_arr))
    temp_df['adv_zero'] = temp_arr
    
    np_path = save_path + '/bertopic_' + str(start) + '_' + str(finish)+'.npy'
    temp_arr = np.load(np_path)
    print('For bertopic ', str(start), ':',str(finish), temp_df.shape[0],'for', np_path, len(temp_arr))
    assert(temp_df.shape[0] == len(temp_arr))
    temp_df['adv_bertopic'] = temp_arr
    dfs.append(temp_df)

    start = finish
    
clean_df = pd.concat(dfs)
clean_df.to_csv(save_path + '/' + os.path.basename(data_path).split('.')[0] + '_filtered.csv', sep=',')

print('------------------------')
print('Filtration is finished')

