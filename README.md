# stockpile
**Language**: Go  

**Description**:  
A resource collater library, which enables an application to access resource usage i.e. CPU percentage, Memory consumption etc. statistics from different machines via SSH and plot the collated data on to a local spreadsheet for analysis.

## Release Notes
##### Version 1.1
###### New Features
- Add Swap Memory in MB if available

##### Version 1.0
###### Initial Features
- Basic system CPU percentage usage
- Memory Used (Percentage and MB)
- Stored in .csv file
- Each record tagged with the master machine's local timestamp
- Custom fields provided by application layer via callback
