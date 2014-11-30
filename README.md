go-lib
======

golang library

bitmap
------
NUMA CPU bitmap, used to affinity progress to some CPU.
CPU must be hypethreaded, and CPU number look like as follows:
[node0, node1, ... , node0, node1, ...]
For example:
node0: [0,1,2,3,4,5,12,13,14,15,16,17]
node1: [6,7,8,9,10,11,18,19,20,21,22,23]
