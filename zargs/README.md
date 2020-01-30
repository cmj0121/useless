# zargs #
The zargs is the Go-lang based argument parse library.

## Option ##
The **Option** is the interface that used as the zargs's flag or command. The flag in zargs
is used as the extra parameter, like **-f** of **--flag** which may has following value. Also
the flag could be the required which means the flag should used in the arguments. On the other
hand, the command or sub-command is the option that append alone and only save the value in zargs.


Both the option has its unique name in the zargs, and sub-command has dependent namespace which means
sub-command can override the option which already save in parent command. For example the argument
**--flag 1 sub --flag 2**, the zargs will get the value **2** for the option **flag** in sub-command **sub**.

