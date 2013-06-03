/*
Package Goodies is a container for several utilities

Sprockets

Integrates the power of sprockets ( http://getsprockets.org ) with your go program

I18N

A gettext based translation library with XLiff support

Ratelimiter

Implements an HTTP rate limiter (including net/http or gorilla) to prevent bad clients from overloading your site.
In-memory and memcached backends currently supported.

Mailer

Provides email delivery (current delivery methods are 'postmark' and 'test')

Monitor

Records how long a function takes to execute

Monitor - db

Wraps an arbitrary database driver, logging all queries made with durations

Monitor - statsd

Reports timings to statsdb; includes an HTTP decorator for logging web transaction timings

Errorhandling

An HTTP decorator which recovers from panics in your server

Goanna

A web toolkit for go; routing, controllers, sessions, views etc.

HTTP - cachecontrol

An HTTP decorator which adds cache-control headers to outgoing responses

HTTP - Debug

An HTTP decorator which dumps request / response bodies to STDOUT

HTTP - Log

An HTTP decorator which logs all requests in Common Log Format

Config

Parse JSON-formatted config files.

Expvar

Imports the 'expvar' builtin package, and exports a couple of useful flags

*/
package goodies
