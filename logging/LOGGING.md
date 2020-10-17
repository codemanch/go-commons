# go-commons logging
The go-commons ```logging``` package provides enterprise grade logging capabilities.

# Features
* Multiple logging levels ```OFF,ERROR,INFO,DEBUG,TRACE```
* Console and File based writers with rolling file support
* Ability to specify log levels for a specific package  
* Internationalisation (i18n) support
* Async logging support
* Configurable can be done using either a file,env variables,code defintions.

## Usage

### Simple Usage 
The simplest example is  as shown below.
```
    package main
    
    import (
        "github.com/appmanch/go-commons/logging"
    )
    
    //logger Package Level Logger
    var logger = logging.GetLogger()
    
    func main() {
        logger.Info("This is an info log msg")
        logger.Warn("This is an warning msg")
        logger.Error("This is an error log msg")
        logger.Error("This is an error log msg with optional error",errors.New("Some error happened here "))
        logger.Warn("Message with any level can also have an error",errors.New("Some error happened here but want it as info"))
    
    }
    
   ```
This will create a logger with default configuration

* Logs get written to console
* Log levels  ```ERROR, WARN``` get written to stderr and remaining levels to stdout
* The default log level is ```INFO``` this can be overwritten using an env variable or log config file. 
  See Log [Configuration](#Log Configuration) section for more details.

###




# Log Configuration
The below table specifies the configuration parameters for logging
The log can be configured in the following ways.

### 1. File Based Configuration
The file based configuration allows a file with log configuration to be specified. Here is a sample file configuration
based on ```json```. 
```
{
     "format": "json",
     "async": false,
     "defaultLvl": "INFO",
     "includeFunction": true,
     "includeLineNum": true,
     "pkgConfigs": [
       {
         "pkgName": "main",
         "level": "INFO"
       }
     ],
     "writers": [
       {
         "console": {
           "errToStdOut": false,
           "warnToStdOut": false
         }
       }
       {
         "file": {
           "defaultPath": "/tmp/default.log"
         }
       },
       
     ]
   }
```

following table specifies the field values
|Field Name   | Type    | Description   | Default Value|
|:-|:-|:-|:-:|
|format|String| The output format of the log message. The valid values are `text` or `json`| `text` |
|async|Boolean| Determines if the message is to be written to the destination asynchronously. If set to `true` then the LogMessage is prepared synchronously.However, it is written to destination in a async fashion.|`false`||defaultLvl| String|Sets the default Logging level for the Logger. This is a global value. For overriding the log levels for a specific packages use `pkgConfigs`. The valid values are `OFF,ERROR,INFO,DEBUG,TRACE`| `INFO`|
|includeFunction| Boolean| Determines if the Function Name needs to be printed in logs. |`false`|
|includeLineNum|Boolean| Determines if the line number needs to be printed in logs. This config takes into effect only if `includeFunction=true`|`false`|
|pkgConfigs   |Array|This field consists array of package specific configuration.<br>`{"pkgName": "<packageName>","level": "<Level>"}`|`null`|
|writers| Array|Array of writers either `file` or `console` based writer. Console Writer Has the following has the following configuration `{"console": {"errToStdOut": false,"warnToStdOut": false}}`.Log levels except `ERROR and WARN` are written to `os.Stdout`. The Entries for the `ERROR,WARN` can be written to either os.StdErr or os.Stdout <br> For a file based log destination,paths for each level can be specified as follows.<br> `{"file": { "defaultPath": "<file Path>","errorPath": "<file Path>","warnPath": "<file Path>","infoPath": "<file Path>","debugPath": "<file Path>","tracePath": "<file Path>" }`<br> If any of the `errorPath,warnPath,infoPath,debugPath,tracePath` is not specified then default path for that level is applied. If all of the level specific paths are specified then the `defaultPath` value is ignored.| N/A|

### 2. Environment Based Configuration
