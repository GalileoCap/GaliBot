import pandas as pd
import os

from utils import *
import cfg

def dfFpath(name):
  return f'{forcePath(cfg.dbDir)}/{name}.csv.gz';

def getDf(name):
  if os.path.exists(dfFpath(name)): return pd.read_csv(dfFpath(name), index_col = 0); 
  else: return pd.DataFrame();

def saveDf(df, name):
  df.to_csv(dfFpath(name));
