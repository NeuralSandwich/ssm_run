# SSM Run

## Description

Tool to run CLI commnds on instances using SSM.

Does require the instance you are looking to access to be associated with SSM.

## Usage

Usage instructions with code examples

```shell
# Here is the code example
ssm_run help                                               # Get Help
ssm_run run -i <instance_id> --command-string "ipconfig"   # Run ipconfig on windows instance
```

## TODO

- [x] Update README
- [ ] Add status command
- [ ] Add automated check to change SSM Document based on instance platform type
