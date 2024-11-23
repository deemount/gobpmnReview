# To Do's

1. Get rid of strings.Contains (it was at first a good entry level to fullfill my ideas, but it is not safe).
2. Define the ruleset for BPMNType, Anonym and Config in an .md file.
3. Describe, how the ProcessName for the single and multiple process is generated.
4. Describe, for what the ProcessName is needed and where it's works.
5. Describe, how extractPrefixBeforeProcess works. Look also in countFieldsInProcess.
6. Describe, what v.Target.FieldByName really is in the algorithm. Look countFieldsInProcess in quantity.go
7. Describe, when an Participant is detected.
8. Explain, which methods from the reflect package were used. Make a list of it and with some Output examples.
9. Of course, a diagram can implemented very fast. Let's do it!
10. Don't try to build a bpmn-file with a diagram - you have it already. The day x is in the making!
11. Define another ruleset for a mix of both given examples (and why it should not work like this).
12. So, the ruleset on point 11 is defined (or later then) - go and code a solution for it.
13. Try to find a solution not to count here and then. It would be better, if the code counts everythingAdd the in one rush.
14. Add the incoming and outgoing logics to the elements.
15. Add the SourceRef and TargetRef to the sequence flows.
16. Add sequence flows for a single process
17. Refactor the call method for a single process. look at internal.go
18. Describe all fields in the reflectValue data structure.
19. Explain, where v.TargetNumField is need.
20. Describe, why "From" is an important prefix in the field name.
21. Is "To" in a name of the field a good solution to handle gateways?
22. Implment a counter for edges and shapes., whe the diagram is implemented as a model.
23. Describe, why a multiple process ranges through anonym fields.
24. Describe, why asingle process is ranging through BPMNType.
25. A big field (in my opinion) is the hash distribution algorithm. Actually the algorithm is scalled to a minimum of use case and could grew by the elements given in a process.
26. Explain, what the hash distribution algorithm could mean in gobpmn.
27. Describe the "typ"-function.

Note: This list will grow in the next days.
