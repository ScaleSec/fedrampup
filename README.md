# FedRAMPup!

How annoying is it to get an inventory of AWS resources for FedRAMP and then put that into the SSP formatted CSV? Very. `fedrampup` will handle all of this for you. It can be run on it's own from a Docker container or in AWS in Fargate with S3 output.


# Configuration

Everyone's environment is different. That's why `fedrampup` has an extensive configuration framework based in environments in keeping with a [12 factor application](https://12factor.net/).


FedRAMP SSP expects the assets in a CSV with the rows:

Unique Asset Identifier
IPv4 or IPv6 Address
Virtual
Public
DNS Name or URL
NetBIOS Name
MAC Address
Authenticated Scan
Baseline Configuration Name
OS Name and Version
Location
Asset Type
Hardware Make/Model
In Latest Scan
// Software and Database Inventory
Software/Database Vendor
Software/Database Name and Version
Patch Level
Function
// Any Inventor
Comments
Serial Number/Asset Tag Number
VLAN/Network ID
System Administrator/Owner
Application Administrator/Owner
