# gobpmn - gelingt mir hier ein erstes, echtes Review?

## Die Idee

Stell dir mal vor, man kann, wie in meinen Beispielen, funktionsfähige Prozesse in Go darstellen und später zur Ausführung in eine BPMN-Engine weiterleiten. Das war es dann auch schon mit der Idee. So einfach kann das sein. Spart einen zusätzlichen Modeller und so manch Go-Entwickler bleibt voll und ganz in seinem Habitat.

## Grundlegendes

### Wie ich ans Werk ging

Ich bediene mich der Extreme-Programming-Methodologie, d.h. keine Tricks, sondern einfaches herunterspulen von Variablen, Schleifen und anschließender Refaktorisierung. Klassischer Aufbau für Anfänger rundet das Gesamtbild, einer sich in Entwicklung befindlicher Anwendung, ab. Ich glaube, sonst wäre ich bis Heute noch längst nicht soweit. Außerdem wirst du, wie Go halt so tickt, viel Pass-By-Value in meinen aufgebauten Methodiken finden. Und manches ist sicherlich noch nicht sehr schön, aber hey, es funktioniert.

Das Projekt diente zuerste meiner eigenen Weiterbildung in Go. Ich mag Go sehr. Angefangen habe ich mit PHP anno 1999 in der dritten Version.
Ich bin glücklich, mich nicht mehr mit PHP im privaten zu beschäftigen. Nun ist aus dem Projekt etwas entstanden, bei dem ich glaube, der Go-Community etwas zurückgeben zu können. Aber das müssen andere (für sich) entscheiden. Vielleicht auch Du?

Den jetzigen Stand der Idee würde ich als dritte Version betiteln. Vorherige Versionen sind auf github, in meinem Account, zu finden.
Die zweite Version, bei der ich eine Aufteilung, einer vorher einzigen Repository, vorgenommen habe, befindet sich unter https://github.com/deemount/gobpmn. Die allererste Version befindet sich in https://github.com/deemount/gobpmnlab.

***Um weitere Kopfschmerzen zu vermeiden, bitte ich, sich zuerst auf dem Stand der dritten Version zu konzentrieren.***

Die erste Version meines Ansatzes unterscheidet sich auch noch dahingehend zu den anderen Versionen, dass ein Prozess in BPMN als vollständiges Modell mit Diagram dargestellt werden kann. Ist aber sehr aufwendig konzipiert. Die nachfolgenden Versionen beschäftigen sich nicht mit einem zusätzlichen Diagram, sondern nur mit der Darstellung eines Prozesses im BPMN-Format von Camunda 7.

Der jetzige Stand ist von den Elemente und deren Felder, im Gegensatz zu den (fast) vollständigen Datenstrukturen in den anderen Versionen, stark gekürzt. Ich musste das machen, weil ich a) das Ganze irgendwie noch auf Medium.com wiedergeben will und um b) schneller weiterzugelangen, ohne mich mit zusätzlichen "Datenmüll" (der ja keiner ist) zu beschäftigen.

### Noch zwei Werkzeuge

Ich nutze den Camunda Modeller, der auch die Version 7 darstellen kann, um die Aufteilung eines Business Process Model in XML zu verstehen und nachzubauen. Ich rate daher, sich dieses Tool selbst zu installieren, um besser zu verstehen. Außerdem empfehle ich dir noch von der OMG, die Dokumentation zu BPMN herunterzuladen.

### Der Ansatz als Blaupause und Datenstruktur

Der Ansatz befindet sich, als zwei verschiedene Prozesse in structs getarnt, in der main.go-Datei wieder.
Einmal einen Beispielprozess und einen Prozess, der eine Vermietung darstellen kann (aber nicht tut).

Der erste Ansatz, als ExampleProcess, gekennzeichnet, ist als BPMN-Datei im Verzeichnis "blueprint" zu finden.
Dies gilt auch für den zweiten Ansatz. Beide Dateien sind nach ihren Beispielen aus der main.go benannt.

### Notizen in den Dateien

Ich habe meine Probleme und Erkenntnisse als sogenannte Notes mit Doppelpunkt gekennzeichnet. Das erspart mir hier eine allzu aufwendige Dokumentation.

### Keine Fehlerbehandlung

Na ja, so gut wie keine. Ich mach das immer lieber zuletzt. Meine Fehler kenn ich ja nun selbst und wenn einer unbehandelt auftaucht, dann siehe Notes bzw. auch als Bug gekennzeichnet (und einigermaßen beschrieben). Ich befinde mich in der Entwicklung und habe keinen Druck.

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

**Hinweis:** Ich nutze das funktionelle Option-Pattern für mehrere Argumente im BPMN-Modeller.

## Ablauflogik in der reflect_di.go

In dieser Datei habe ich gerade nun gar keine Lust, einen Programmablaufplan (PAP) zu erstellen, noch eine Aufzählung der verschiedenen Ausführungen zu beschreiben. Stattdessen erzähle ich hier den Ablauf kurz und bündig als Geschichte.

### Die Geschichte

Ich habe eine reflektierte Instanz in v.Target, in der ich zuerst einmal die Prozesse "auffange" und mit Hilfe von Reflektionen mit Werten, wie den Typen und Hash "ausstatte" (das nennt sich dann wohl reflektierte Abhängigkeiteninjektion). v.Target überträgt dann die Werte in v.Def (auch wieder reflektierte Abhängigkeiteninjektion). Dazwischen gibt es einige Methoden, die ich versucht habe, den Aufgaben nach, einem Namen zu geben (mehr oder weniger gut gelungen), die diese Schritte dann ausführen.

## Ablauflogiken in quantity.go und mapping.go

In meinen bisherigen Ansatz muss ich viel zählen. Ich muss wissen, wieviele Elemente ich erstellen muss. Die Schwierigkeit besteht hier auf jeden Fall darin, eine richtige Indexierung anhand von Schlüsselwerten vorzunehmen. Dazu hilft mir die quantity.go-Datei.

Um zu unterscheiden, mit was für einen Prozessaufbau ich es eigentlich zu tun habe, benötige ich eine Aufteilung, die ich in der mapping.go vornehme. In meinem Regelsatz befindet sich auf jeden Fall ein Versuch, eine Kollaboration mittels anonymer Felder in einer Datenstruktur zu erkennen und dann auszuführen.

## Die internal.go

Hier befinden sich ausgelagerte Methoden, die ich als Funktionen nutze. Einfach, um nicht bis auf Zeile 1000 herunterscrollen zu müssen. Ich habe die Abschnitte für die jeweilige Ausführung eines Beispieles gekennzeichnet, so dass man nicht den Überblick verliert. Auch hier befinden sich sogenannte Notizen und Beschreibungen.

## Wünsche für die Zukunft

Also ein Diagram muss noch dazu kommen. Sicherlich! Aber dazu brauch es dann auch noch ein wenig mehr Kenntnisse. Meine ersten Gehversuche habe ich dahingehend in der Repository gobmnLab (siehe Verzeichnis "examples" oder "models/bpmn/canvas") gesammelt. Unter anderem musste ich mich mit Fixed-Point-Notation und Linienüberschreitung beschäftigen. Auch die Ermittlung der Wegpunkte ist noch einmal ein Thema für sich. Alles in Allem ist eine Kompilierung eins Diagrams mittels Reflektion weitaus aufwendiger, als nur die Darstellung eines Prozesses selbst (Erkenntnisgewinn).

Auch bei der Ausgabe der verschiedenen Formate, wünsche ich mir in Zukunft die Möglichkeit, die selbe Reihenfolge darzustellen, wie sie in den Beispieldatenstrukturen angeordnet sind. Im Moment wird die Reihenfolge dargestellt, wie sie in der elements.go, in der Prozess-Datenstruktur (Proces) zu finden ist.

Was mir aber im Moment mit am Wichtigsten ist, andere Entwickler für diesen Ansatz zu begeistern und Unterstützung aus freien Willen zu erfahren. Da draußen muss es doch noch irgendwie noch Genies geben, die das Ganze als hilfreich empfinden und meinen Ansatz eventuell auffassen und noch besser machen können.
