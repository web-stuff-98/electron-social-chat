# UNFINISHED group video/voip chat app using Vue, Tauri & Go

## not using electron anymore, shouldn't be called "electron" social chat

## The serverside code is pretty much the same as my last project, except sessions are kept on Redis, where they should be.

Models for socket events are different. The client no longer sends messages to request to open subscriptions directly, so there are many more socket models but its easier to understand what is going on

It was using Electron beforehand, but I changed to tauri because its much faster and lightweight

SCSS classes are kebab cased this time, and I'm using functions instead of arrow functions for a lot of stuff