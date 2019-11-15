# Jeeves

A discord bot for your Classic WoW guild.

## Introduction

Every guild has a Discord channel that acts as the hub for the community. Unforunately,
Discord is missing a lot of features that any self-respecting guild requires to function.
Jeeves understands that you don't want to split your guild's attention between Discord
and a clunky website. To help with this, he provides a number of different services
for your guild including:

- Guild Bank

## Guild Bank

### Setup

Before Jeeves can track your guild bank's contents, you must tell him where he can display what is
inside of your bank. To do this, an admin has to say `!jeeves-assign-bank` inside of the desired channel.

Note: Jeeves will also create the `Banker` role if your discord doesn't have it. You can assign this role to users
that you want to be able to move items in and out of the bank, without elevating them to an Admin.

#### Depositing Items

Depositing an item is done by sending a message of the form:

```
!deposit Lava Core
```

You can deposit multiple items:

```
!deposit Lava Core, Righteous Orb
```

As well as multiple quantities of items (specifying `1x` is optional):

```
!deposit 3x Lava Core, 2xRighteous Orb
```

NOTE: item names are not case-sensitive but the command is, so `!deposit lava core` would have also worked.

#### Withdrawing Items

Withdrawing items looks very similar to depositing:

```
!withdraw Lava Core
```

You can withdraw multiple items:

```
!withdraw Lava Core, Righteous Orb
```

As well as multiple quantities of items (specifying `1x` is optional):

```
!withdraw 3x Lava Core, 2xRighteous Orb
```

NOTE: item names are not case-sensitive but the command is, so `!withdraw lava core` would have also worked.
