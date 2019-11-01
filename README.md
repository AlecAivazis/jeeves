# Jeeves

A discord bot for your Classic WoW guild.

NOTE: none of this is implemented... yet

## Introduction

Every guild has a Discord channel that acts as the hub for the community. Unforunately,
Discord is missing a lot of features that any self-respecting guild requires to function.
Jeeves understands that you don't want to split your guild's attention between Discord
and a clunky website. To help with this, he provides a number of different services
for your guild including:

- Guild Bank

## Setup

Jeeves works across many channels in your guild. Once you have invited Jeeves to your guild, you
must assign the neccessary channels for each feature to work. In order to assign a channel, navigate
to the appropriate channel and type `!jeeves-assign-channel <role>`.

| Channel Role | Usage                                                                                    |
| ------------ | ---------------------------------------------------------------------------------------- |
| bank         | The channel that will have the summary of the bank                                       |
| bank-log     | The channel that the bank manager will use to record items moving in and out of the bank |

## Roles

Most guilds want to limit restrict certain actions like managing the guild bank to a subset of users.
Jeeves relies on Discord's Role system to enforce these restrictions. When Jeeves is first added to your
guild, it will create the necessary roles.

| Role Name | Permissions                                           |
| --------- | ----------------------------------------------------- |
| Banker    | Allowed to record items moving in and out of the bank |

## Guild Bank

Jeeves will track items as they move in and out of your guild's various bank alts. At the moment, there is
no in-game addon to aid this process. If something like that would interest you, please open an issue!

Once you have assigned the `bank` and `bank-log` roles,
