# Logray

## Goals

### Multiple outputs 

Our previous library had this, and logrus does through hooks. This is important to start using syslog. 

(As an aside, the hooks approach is odd. In terms of practicality, you could view writing to syslog as more of an output instead of hook. Additionally, it makes sense to configure it as you would an output, with things like the loglevels it should send there. Additionally, viewing the formatter and output as two separate things like logrus does does not mesh with our outstanding).

### Field propagation

Logrus added a fields map which allowed for more valuable metadata. We liked this and wanted to extend it, but logrus didn’t have fields building on fields. We wanted to take a context and pass it further down into the system, and have the fields build upon it as its actions got more specific. Logrus only lets you do a WithFields, which gave you an Entry, and could write multiple things on that, but couldn’t really have it turtle further down.

### Contextually specific log formatting

At different places in the system, you care about different data related to log messages. The logging changes allow formats/outputs to be overwritten in certain places where the format we see directly matters in a slightly different way.

### Simple to use

This library needed to be simple to use. So starting out, basically all you have have to do is write logger := logray.New() to start logging. (Note, perhaps writing to stdout by default should be added to the library). Want to add fields? You can write logger.SetField(key, value). Want to pass it further down where separate fields might be set? That's just logger2 := logger.Clone().
