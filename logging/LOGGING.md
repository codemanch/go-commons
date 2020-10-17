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
| Field Name   | Description |
|  :---:       | ------------- |
| format | Either text or Json  |
| | Content Cell  |

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
         "file": {
           "defaultPath": "/tmp/default.log"
         }
       },
       {
         "console": {
           "errToStdOut": false,
           "warnToStdOut": false
         }
       }
     ]
   }
```

following table specifies the field values
|Field Name   | Type    | Description   | Default Value   |
|:-:|:-:|:-|:-:|
|format|String| The output format of the log message. The valid values are `text` or `json`| text |
|async|Boolean| Determines if the message is to be written to the destination asynchronously. If set to `true` then the LogMessage is prepared synchronously and then the write at destination in a async fassion.    |`false`|
|   |   |   |   |