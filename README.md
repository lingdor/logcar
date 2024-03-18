# logcar
a log middleware to direct streams, rolling, filter, split to files  kafka, hdfs or others.

case 1:
```text
|--progress---|                |-----------------|
|  app        | ---pipline---> |    logcar       |
|-------------|                |-----------------|
```
case 2:
```text

           Dev  /  Ops
|--container--|    |----------|     |-----------------|    |----------------------|
|  app        | -> |  logcar  | ->  |  kafka          | -> | logcar file-appender |
|(go|java|...)|    |----------|     |-----------------|    |----------------------|
|-------------|                            |               |----------------------|
                                           |-------------> |        flink/ES      |
                                                           |----------------------|
```
case 3:

```text

svc1                        svc3
[ pod 1] -----------↓       [pod logcar]
[ pod 2] --------------->   [pod logcar] --> svc1.log 
[ pod 3] -----------↑               ↑    --> svc2.log
svc2                                |
[ pod 1] ---------------------------|
[ pod 2] ---------------------------|
[ pod 3] ---------------------------|



```


```
# Command

```shell
./xx-server|logcar -f ./logcar.yml
```




