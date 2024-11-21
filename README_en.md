# gobpmn - do I manage here a first, real review?

## The idea

Imagine you can represent functional processes in Go, as in my examples, and later forward them to a BPMN engine for execution. That's all to the idea. It can be that simple. Saves an additional modeler and many Go developer remains entirely in their habitat.

## Fundamentals

### How I got to work

I use the extreme programming methodology, i.e. no tricks, but simple rewinding of variables, loops and subsequent refactoring. A classic structure for beginners rounds off the overall picture of an application under development. I don't think I would have gotten that far yet. Also, as Go ticks, you will find a lot of pass-by-value in the methodologies I have developed. And some of it is certainly not very pretty yet, but hey, it works.
The project was first of all for my own further training in Go. I really like Go. I started with PHP in 1999 in the third version.
I am happy that I no longer have to deal with PHP in my private life. Now the project has turned into something where I believe I can give something back to the Go community. But that's for others to decide (for themselves). Maybe you too?
I would call the current status of the idea the third version. Previous versions can be found on github, in my account.
The second version, in which I have split a previously single repository, can be found at [gobpmn](https://github.com/deemount/gobpmn). The very first version can be found at [gobpmnLab](https://github.com/deemount/gobpmnlab).

***To avoid further headaches, please concentrate on the status of the third version first***

The first version of my approach also differs from the other versions in that a process can be represented in BPMN as a complete model with a diagram. However, this is a very complex concept. The following versions do not deal with an additional diagram, but only with the representation of a process in the BPMN format of Camunda 7.
The current version, in contrast to the (almost) complete data structures in the other versions, is very abbreviated in terms of the elements and their fields. I had to do this because a) I want to somehow reproduce the whole thing on Medium.com and b) to get on faster without having to deal with additional “data garbage” (which is not really garbage at all).

### Two more tools

I use the Camunda Modeller, which can also display version 7, to understand and recreate the breakdown of a business process model in XML. I therefore advise you to install this tool yourself in order to understand it better. I also recommend that you download the documentation on BPMN from the OMG.

### The approach as a blueprint and data structure

The approach is disguised as two different processes in structs in the main.go file.
One is an example process and the other is a process that can (but does not) represent a rental.
The first approach, labeled as ExampleProcess, can be found as a BPMN file in the “blueprint” directory.
This also applies to the second approach. Both files are named after their examples from main.go.

### Notes in the files

I have marked my problems and findings as so-called notes with a colon. This saves me from having to write too much documentation here.

### No Error Handling

Well, almost none. I always prefer to do it last. I know my bugs myself and if one appears untreated, then see Notes or marked as a bug (and described to some extent). I'm in the development phase and have no pressure.

## Sequence logics in the main.go

I have been too lazy to create a program flow chart so far. Therefore, I will present the flow logic in this file as an enumeration.

1. initialize the method NewReflectValue from the file “reflect_di.go” in the main function.
2. initialize the method NewBPMNModeller from the file “modeller.go” in the main function.
3. subsequently call the method Marshal in the main function.
4. in the directory “files” there are second directories for the respective file formats, json and bpmn.
   The created processes are saved here.

These four steps are executed when you use the following command in your terminal:

```go
go run *go
```

**Note:** I use the functional option pattern for multiple arguments in the BPMN modeler.

## Sequence logics in the reflect_di.go

In this file, I don't feel like creating a program flow chart or describing a list of the various versions. Instead, I will tell you the process in a short and concise story.

### The Story

I have a reflected instance in v.Target, in which I first “catch” the processes and “equip” them with values such as the types and hash with the help of reflections (this is probably called reflected dependency injection). v.Target then transfers the values to v.Def (again reflected dependency injection). In between, there are some methods that I have tried to name according to the tasks (more or less successful), which then perform these steps.

## Sequence logics in quantity.go and mapping.go

In my previous approach, I have to count a lot. I need to know how many elements I need to create. In any case, the difficulty here is to make a correct indexing based on key values. The quantity.go file helps me with this.
In order to distinguish what kind of process structure I am actually dealing with, I need a breakdown, which I make in the mapping.go file. In my rule set, there is definitely an attempt to recognize and then execute a collaboration using anonymous fields in a data structure.

### The internal.go

Here you will find outsourced methods that I use as functions. Simply to avoid having to scroll down to line 1000. I have labeled the sections for the respective execution of an example so that you don't lose track. There are also so-called notes and descriptions here.

## Wishes for the future

So a diagram still has to be added. Certainly! But that requires a little more knowledge. I made my first attempts in the gobmnLab repository (see directory “examples” or “models/bpmn/canvas”). Among other things, I had to deal with fixed-point notation and line crossing. Determining the waypoints is another topic in itself. All in all, compiling a diagram using reflection is far more complex than just displaying a process itself (gaining knowledge).

In the future, I would also like to be able to display the different formats in the same order as they are arranged in the sample data structures. At the moment, the sequence is displayed as it can be found in elements.go, in the process data structure (Process).

But what is most important to me at the moment is to inspire other developers for this approach and to receive support of their own free will. There must still be some geniuses out there who find the whole thing helpful and who might be able to grasp my approach and make it even better.
