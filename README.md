Vole
====

Vole is a web application for sharing words, pictures and videos with others. Data is distributed peer-to-peer using Bittorrent Sync.

A complete introduction is available on [vole.cc](http://vole.cc).

If you're planning on using Vole, please join the [mailing list](https://groups.google.com/forum/#!forum/vole-dev) for announcements and development discussion.

Getting started
---------------

Currently we are pre-release. The following steps are for interested developers.

* [Download Vole](https://www.dropbox.com/sh/cd943dyjjavzxk7/uqm3wh37fS) for Mac OSX, Windows or Linux. The current version is 0.1.0.
* Or, you can compile it yourself. See [CONTRIBUTING](https://github.com/vole/vole/blob/master/CONTRIBUTING.md).
* Run the Vole application from the command line. First, `cd` into the directory that you placed Vole. Then, on Mac/Linux type `./vole`, and on Windows type `vole.exe`.
* Open a web browser and go to [http://localhost:6789](http://localhost:6789).
* Click 'My Profile'
* Enter your name. Enter your Gravatar email (optional).
* Click 'Home'
* Post something.

Sharing and following
---------------------

Start by installing [Bittorrent Sync](http://labs.bittorrent.com/experiments/sync.html).

Following and sharing currently involves manually setting up folders, however we'd like to automate it as soon as Bittorrent Sync releases a build that supports control via an API.

Following others
----------------

* Grab the **read-only** ID of the person you want to follow. A directory is in progress at [vole.cc](http://vole.cc). Why not start with Vole Team updates? Our key is RA32XLBBHXMWMECGJAJSJMMPQ3Z2ZGR7K.
* Find the Vole `users` folder. Unless you changed the defaults, it will be in a directory called `Vole/users` in your home folder.
* Create a new folder in `Vole/users`, you should name it after the user that you're about to follow. For example, `Vole/users/voleteam`.
* In BitTorrent Sync, add this new folder as a shared folder, using the read-only key you grabbed in step 1. [Instructions](http://labs.bittorrent.com/experiments/sync/get-started.html) are available on their site and vary a little by operating system.
![OSX Screenshot](https://f.cloud.github.com/assets/453297/692312/c113737a-dc18-11e2-84e4-dee7e0507c08.png)
* You should receive notification that the folder has sync'd.
* In your browser, see the new posts appear.

Sharing your posts
------------------

Find your own user folder, for example, if you created a profile named 'Chuck':

    /home/chuck/Vole/users/Chuck_9674e8e8-7c7a-41e6-52ed-51a3f7969812

* In Bittorrent Sync, add this folder as a shared folder.
* In the folder options, grab the **read-only key**. Make sure the key starts with the letter 'R' that signifies it's the read-only one. You can find it by going to the advanced folder preferences. This is the key that you can share with others so they can follow your posts.

Profile site
------------

We plan on allowing people to claim their names on vole.cc and store their keys along with their profile, making it much easier to find people who are posting publicly.

Versions
--------

Vole uses [semantic versioning](http://semver.org)

Technology
----------

* [Bittorrent Sync](http://labs.bittorrent.com/experiments/sync.html)
* [Go](http://golang.org/)
* [Ember.js](http://emberjs.com/)

License
-------

Copyright (C) 2013 Vole development team

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
