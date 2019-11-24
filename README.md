# Jeeves

A discord bot for your Classic WoW guild.

## Introduction

Every guild has a Discord channel that acts as the hub for the community. Unforunately,
Discord is missing a lot of features that any self-respecting guild requires to function.
Jeeves understands that you don't want to split your guild's attention between Discord
and a clunky website. To help with this, he provides a number of different services
for your guild including:

- Guild Bank
- Calendar and Event Scheduling (coming soon)
- Profession/Recipe look up (coming soon)

## Guild Bank

### Setup

Before Jeeves can track your guild bank's contents, you must tell him where he can display what is
inside of your bank. To do this, an admin has to say `!jeeves-assign-bank` inside of the desired channel.

Note: Jeeves will also create the `Banker` role if your discord doesn't have it. You can assign this role to users
that you want to be able to move items in and out of the bank, without elevating them to an `Admin`.

### Commands

A `Banker` can interact with Jeeves through a few different commands summarized below. Note that item names are
not case-sensitive but the command name is. So, `!deposit lava core` works but `!DEPOSIT Lava Core` does not work.

| Command  | Description                       | Example                                             |
| -------- | --------------------------------- | --------------------------------------------------- |
| deposit  | Adds items to the guild bank      | `!deposit Lava Core, 2xRighteous Orb, 2g, 30s, 15c` |
| withdraw | Removes items from the guild bank | `!withdraw Lava Core, 2x Righteous Orb`             |
