# Packer Fantasy #
The packer fantasy (PF) is a simple but useless binary packer based on Go-lang.
This library support simple binary packer functionality, include compression, polymorphic
and obfuscation.


## Code Architecture ##
In the PF, there are two major parts of code: binary tool and packer. The binary tool
is used to analysis the source file and abstract the parts of machine code, and the
packer is used to compress and obfuscate and generate the final code.
