# gostip        [![English](gb40.png "Union Jack")](#en) [![Deutsch](de33.png "Deutsche Flagge")](#de)

##### <a id="en"></a> Go web application for school graduates (kyrgyz: abiturients) who want to apply for a DAAD grant at KGFAI

Besides its primary purpose, this application is used as part of a tutorial to demonstrate use of different libraries 
to write a complete and working web application.

##### <a id="de"></a> Go Web Anwendung für Abiturienten, die sich für ein DAAD Stipendium an der DKFAI bewerben wollen
###### Installation
go get lädt nicht die transitiven Abhängigkeiten, daher werden diese mit go get ... aufgelöst. Dabei gibt es einige 
harmlose Fehlerausgaben.
```
go get -u github.com/geobe/gostip
go get ...
```
Außerdem wird postgresql verwendet. Dazu muss

* [pgSql](https://www.postgresql.org/download/) installiert und
* eine neue Datenbank mit einem Benutzer für die Go Applikation angelegt werden. 
Das geht einfach mit [pgAdmin](https://www.pgadmin.org/).
* Für [pgAdmin und pgSql](http://www.enterprisedb.com/products-services-training/pgdownload) gibt es einen gemeinsamen 
Installer von EDB.
* Datenbankname, Username und Passwort werden in der 
[Konfigurationsdatei](https://github.com/geobe/gostip/blob/master/config/devconfig.json) festgelegt.
* In der pgSql Konfigurationsdatei postgresql.conf timezone auf 'UTC' setzen. Mit anderen timezone Einstellungen gibt es 
möglicherweise Fehler.
 
######Entwicklung
Die Applikation wird mit IntelliJ Community Edition entwickelt. Sinnvollerweise legt man dazu ein 
Go Projekt im $GOPATH Verzeichnis an. Dann kann das Projekt aus der IDE getartet werden.

######Tutorial
Neben dem eigentlichen Zweck der Anwendung ist dies gleichzeitig eine Demo-Applikation für ein 
go Webapp Tutorial, in der verschiedene Bibliotheken benutzt werden, um eine
vollständige Webanwendung zu entwickeln:
*   jinzhu/gorm für die persistente Speicherung von Go Objekten
*   gorilla web framework Komponenten für das Multiplexen von Http Requests, Sitzungsverwaltung, Login, xsrf protection und vieles mehr
*   justinas/alice für handler chaining
*   dchest/captcha für Captchas (completely automated public Turing test to tell computers and humans apart)
*   spf13/viper für Konfigurationsfiles
*   Bootstrap3 für responsives Web-Design
*   jQuery für dynamische Webseiten
