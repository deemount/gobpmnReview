# gobpmn - erster Versuch eines Review

## Grundlegendes

Ich bediene mich der Extreme-Programming-Methodologie, d.h. keine Tricks, sondern einfaches herunterspulen von Variablen, Schleifen und anschließender Refaktorisierung. Ich glaube, sonst wäre ich bis Heute noch längst nicht soweit. 

Das Projekt diente zuerste meiner eigenen Weiterbildung in Go. Ich mag Go sehr. Angefangen habe ich mit PHP anno 1999 in der dritten Version.
Ich bin glücklich, mich nicht mehr mit PHP im privaten zu beschäftigen. Nun ist aus dem Projekt etwas entstanden, bei dem ich glaube, der Go-Community etwas zurückgeben zu können. Aber das müssen andere (für sich) entscheiden.

Den jetzigen Stand der Idee würde ich als dritte Version betiteln. Vorherige Versionen sind auf github, in meinem Account, zu finden.
Die zweite Version, bei der ich eine Aufteilung, einer vorher einzigen Repository, vorgenommen habe, befindet sich unter https://github.com/deemount/gobpmn. Die allererste Version befindet sich in https://github.com/deemount/gobpmnlab.

**Um weitere Kopfschmerzen zu vermeiden, bitte ich, sich zuerst auf dem Stand der dritten Version zu konzentrieren.**

Die erste Version meines Ansatzes unterscheidet sich auch noch dahingehend zu den anderen Versionen, dass ein Prozess in BPMN als vollständiges Modell mit Diagram dargestellt werden kann. Ist aber sehr aufwendig konzipiert. Die nachfolgenden Versionen beschäftigen sich nicht mit einem zusätzlichen Diagram, sondern nur mit der Darstellung eines Prozesses im BPMN-Format von Camunda 7.

Ich nutze den Camunda Modeller, der auch die Version 7 darstellen kann, um die Aufteilung eines Business Process Model in XML zu verstehen und nachzubauen. Ich rate daher, sich dieses Tool selbst zu installieren.

### Ansatz

Der Ansatz befindet sich, als zwei verschiedene Prozesse in structs getarnt, in der main.go-Datei wieder.
Einmal einen Beispielprozess und einen Prozess, der eine Vermietung darstellen kann (aber nicht tut).

Der erste Ansatz, als ExampleProcess, gekennzeichnet, ist als BPMN-Datei im Verzeichnis "blueprint" zu finden.
Dies gilt auch für den zweiten Ansatz. Beide Dateien sind nach ihren Beispielen aus der main.go benannt.

### Notizen in den Dateien

Ich habe meine Probleme und Erkenntnisse als sogenannte Notes mit Doppelpunkt gekennzeichnet. Das erspart mir hier eine allzu aufwendige Dokumentation.

### Keine Fehlerbehandlung

Na ja, so gut wie keine. Ich mach das immer lieber zuletzt. Meine Fehler kenn ich ja nun selbst und wenn einer unbehandelt auftaucht, dann siehe Notes bzw. auch als Bug gekennzeichnet (und einigermaßen beschrieben)

## Ablauflogik in der main.go

Ich war bis jetzt zu faul, einen Programmablaufplan zu erstellen. Daher werde ich die Ablauflogik in dieser Datei als Aufzählung darstellen.

1. Initialisiere in der main-Funktion, die aus der Datei "reflect_di.go" stammende Methode, NewReflectValue.
2. Initialsiere nachfolgend in der main-Funktion, die aus der Datei "modeller.go" stammende Methode, NewBPMNModeller.
3. Rufe nachfolgend die Methode Marshal, in der main-Funktion, auf.
4. Im Verzeichnis "files" befinden sich zweitere Verzeichnisse für die jeweiligen Datei-Formate, json und bpmn.
   Hier werden die erstellten Prozesse gespeichert.

Diese vier Schritte werden ausgeführt, wenn du in deinem Terminal folgenden Befehl verwendest:

```go
go run *go
```

Hinweis: Ich nutze das Option-Pattern für mehrere Argumente in einem BPMN-Modeller.

## Ablauflogik in der reflect_di.go

In dieser Datei habe ich gerade nun gar keine Lust, einen Programmablaufplan zu erstellen noch eine Aufzählung der verschiedenen Ausführen zu beschreiben. Stattdessen erzähle ich hier den Ablauf kurz und bündig als Geschichte.

Die Geschichte:
Ich habe eine reflektierte Instanz in v.Target, in der ich zuerst einmal die Prozesse "auffange" und mit Hilfe von Reflektionen mit Werten, wie den Typen und Hash "ausstatte" (das nennt sich dann wohl Reflected Dependency Injection). v.Target überträgt dann die Werte in v.Def. Dazwischen gibt es einige Methoden, die ich versucht habe, der Aufgabe nach einem Namen zu geben (mehr oder weniger gut gelungen).

## Wünsche für die Zukunft

Also ein Diagram muss noch dazu kommen. Sicherlich! Aber dazu brauch es dann auch noch ein wenig mehr Kenntnisse. Meine ersten Gehversuche habe ich dahingehend in der Repository gobmnLab gesammelt. Unter anderem musste ich mich mit Fixed-Point-Notation und Linienüberschreitung beschäftigen. Auch die Ermittlung der Wegpunkte, ist noch einmal ein Thema für sich. Alles in Allem ist eine Kompilierung eins Diagrams mittels Reflektion weitaus aufwendiger, als nur die Darstellung eines Prozesses.

Auch bei der Ausgabe der verschiedenen Formate, wünsche ich mir in Zukunft die Möglichkeit, die selbe Reihenfolge darzustellen, wie sie in den Datenstrukturen angeordnet sind. 
