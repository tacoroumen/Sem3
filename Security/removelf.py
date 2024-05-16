#!/usr/bin/python3
import os
import pandas as pd

#define paths
in_file = "/home/ansadmin/rest.csv"
out_file = "/home/ansadmin/report.csv"

while True:

	if os.path.exists(in_file):
		data= pd.read_csv(in_file)
		#data
		data1 = data.replace('\n','',regex=True)
		#print(data1)
		data1.to_csv(out_file, index=False)
		os.remove(in_file)
