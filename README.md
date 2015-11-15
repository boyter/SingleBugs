# SingleBugs
A simple single person bug tracker

Singlebugs is a sole developer focused bug tracking tool. It takes its interface interface 
cues from Outlook's three panel design where the leftmost are the projects, the middle the
issues and the right running commentary. It has a fulltext search bar which searches
across all three streams making a very fluid development experience. Since it is designed
to work locally it is also very fast.

Running it should be as simple as compiling the Golang code and assets. There is a script
compile-normal.sh which should do this for you and leave you with an executable to run. To
run you need to simply run the compiled file, however you can change the port it runs on
via command line options such as

./fizzbuzz -port=8080

This has been used for several years while developing https://searchcode.com/ and since it works
100% offline has been a very useful when working in environments where network connectivity is 
spotty such as the train.

Generally you would want to put a compiled version of SingleBugs in a shared
directory solution such as Dropbox. This will sync the data files across all environments where you
might want to track bugs and doubles as a backup solution.

Screenshots

![General view](http://www.boyter.org/wp-content/uploads/2013/09/screenshot1.png)
![Editing an issue](http://www.boyter.org/wp-content/uploads/2013/09/screenshot2.png)
![Creating a new issue](http://www.boyter.org/wp-content/uploads/2013/09/screenshot3.png)
