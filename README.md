# ViewSonic Projector

This library provides functionality to communicate with ViewSonic projectors using the RS-232 protocol as specified in the v1.19 documentation.

## **ViewSonic Projector RS-232 Command Parsing**

This document outlines how to structure and parse command packets for communicating with ViewSonic projectors via the RS-232 protocol, based on the v1.19 specification.

## **3.1 Command Description**

The communication flow is a simple request-reply model between the host (control unit) and the projector. The host sends either a Write command to set a value or a Read command to get a value, and the projector sends a corresponding response packet.

* **Write-Function**: Used to control the projector and change its settings.  
* **Read-Function**: Used to query the projector for its current status or settings.

### **Length Byte Description (LSB, MSB)**

* **LSB (Low Byte)**: The low byte of the length. The length is the total number of bytes in the command starting from BYTE 5 up to, but not including, the final Checksum byte.  
* **MSB (High Byte)**: The high byte of the length. For this protocol, it is typically 0x00.

## **Packet Formats**

### **Write-Function Packet Format**

This format is used to send a command to the projector to perform an action or change a setting.

| Byte Position | Field Name | Example Value | Description |
| :---- | :---- | :---- | :---- |
| BYTE 0 | Command Head | 0x06 | Start of a Write Packet. (Cmd1) |
| BYTE 1 | Command Head | 0x14 | Protocol identifier. |
| BYTE 2 | Command Head | 0x00 | Protocol identifier. |
| BYTE 3 | LSB (Length) | 0x04 | Low byte of the payload length. |
| BYTE 4 | MSB (Length) | 0x00 | High byte of the payload length. |
| BYTE 5 | Command | 0x34 | Main command identifier. |
| BYTE 6 | Payload | Cmd2 | First part of the specific command code. |
| BYTE 7 | Payload | Cmd3 | Second part of the specific command code. |
| BYTE 8 | Payload | data | Data value to be written for the command. |
| BYTE N+8 | Checksum | checksum | The checksum calculated for the packet. |

### **Read-Function Packet Format**

This format is used to request information or a status value from the projector.

| Byte Position | Field Name | Example Value | Description |
| :---- | :---- | :---- | :---- |
| BYTE 0 | Command Head | 0x07 | Start of a Read Packet. (Cmd1) |
| BYTE 1 | Command Head | 0x14 | Protocol identifier. |
| BYTE 2 | Command Head | 0x00 | Protocol identifier. |
| BYTE 3 | LSB (Length) | 0x05 | Low byte of the payload length. |
| BYTE 4 | MSB (Length) | 0x00 | High byte of the payload length. |
| BYTE 5 | Command | 0x34 | Static value for Read commands. |
| BYTE 6 | Command | 0 | Static value for Read commands. |
| BYTE 7 | Command | 0 | Static value for Read commands. |
| BYTE 8 | Payload | Cmd2 | First part of the command code to query. |
| BYTE 9 | Payload | Cmd3 | Second part of the command code to query. |
| BYTE 10 | Checksum | checksum | The checksum calculated for the packet. |

## **Reply Formats**

### **Write Response Packet (ACK)**

This is the acknowledgment packet sent by the projector after a successful Write-Function command.

| Byte Position | Field Name | Value | Description |
| :---- | :---- | :---- | :---- |
| BYTE 0 | Command Head | 0x03 | Start of a Write Response packet. |
| BYTE 1 | Command Head | 0x14 | Protocol identifier. |
| BYTE 2 | Command Head | 0x00 | Protocol identifier. |
| BYTE 3 | Command Head | 0x00 | Static value. |
| BYTE 4 | Command Head | 0x00 | Static value. |
| BYTE 5 | Checksum | 0x14 | The checksum for the response packet. |

### **Read Response Packet (1-byte data)**

This is the response from the projector for a Read-Function command that returns a single byte of data.

| Byte Position | Field Name | Value | Description |
| :---- | :---- | :---- | :---- |
| BYTE 0 | Command Head | 0x05 | Start of a Read Response packet. |
| BYTE 1 | Command Head | 0x14 | Protocol identifier. |
| BYTE 2 | Command Head | 0x00 | Protocol identifier. |
| BYTE 3 | LSB (Length) | LSB | Low byte of the payload length. |
| BYTE 4 | MSB (Length) | MSB | High byte of the payload length. |
| BYTE 5 | Command | 0 | Static value. |
| BYTE 6 | Command | 0 | Static value. |
| BYTE 7 | Payload | data | The single-byte data value being returned. |
| BYTE 8 | Checksum | checksum | The checksum for the response packet. |

### **Read Response Packet (2-byte data)**

This is the response from the projector for a Read-Function command that returns two bytes of data.

| Byte Position | Field Name | Value | Description |
| :---- | :---- | :---- | :---- |
| BYTE 0 | Command Head | 0x05 | Start of a Read Response packet. |
| BYTE 1 | Command Head | 0x14 | Protocol identifier. |
| BYTE 2 | Command Head | 0x00 | Protocol identifier. |
| BYTE 3 | LSB (Length) | LSB | Low byte of the payload length. |
| BYTE 4 | MSB (Length) | MSB | High byte of the payload length. |
| BYTE 5 | Command | 0 | Static value. |
| BYTE 6 | Payload | data (LSB) | The low byte of the two-byte data value being returned. |
| BYTE 7 | Payload | data (MSB) | The high byte of the two-byte data value. |
| BYTE 8 | Checksum | checksum | The checksum for the response packet. |

## **Application Notes**

The following notes from the specification provide additional details for parsing specific responses and other contextual information.

### **1\. Parsing Multi-Byte Data Responses**

Some Read commands return a multi-byte value that must be converted from hexadecimal to decimal. The byte order is **Little-Endian**.

#### **Operating Temperature (Note 1\)**

* **Value Bytes**: BYTE 7 through BYTE 10\.  
* **Formula**: HEX2DEC(ddccbbaa) / 10 where aa is BYTE 7, bb is BYTE 8, etc.  
* **Example**: If the response bytes are 0x29 0x01 0x00 0x00, the value is 0x00000129.  
  * HEX2DEC(0x00000129) \= 297  
  * 297 / 10 \= **29.7°C**

#### **Light Source Usage Time (Note 4\)**

* **Value Bytes**: BYTE 7 through BYTE 10\.  
* **Formula**: HEX2DEC(ddccbbaa) where aa is BYTE 7, bb is BYTE 8, etc.  
* **Example**: If the response bytes are 0xB8 0x0B 0x00 0x00, the value is 0x00000BB8.  
  * HEX2DEC(0x00000BB8) \= **3000 hours**

### **2\. Special Response Formats**

#### **Error Status (Note 3\)**

* For service debugging, a read command for error status returns a 32-byte packet.  
* The response starts with 0x05 0x14 0x00 0x16 0x00 0x00 0x00...  
* The payload contains 20 items detailing various error states.

#### **Function Disabled Response (Note 5\)**

* If a function is disabled or "greyed out" (e.g., Aspect Ratio when no source is active), the projector will return a specific response.  
* **Identifier**: The first byte of the response is 0x00.  
* **Full Response**: 0x00 0x14 0x00 0x00 0x00 0x14

### **3\. Contextual Information**

#### **LAN Control (Note 2\)**

* When controlling the projector over a LAN, use the same command codes but remove the 0x prefix.  
* Communication occurs via **port 4661**.

#### **HDMI Range (Note 6\)**

* **Enhanced (PC Level)**: Represents a full 0-255 signal range.  
* **Normal (Video Level)**: Represents a limited 16-235 signal range.

#### **Projector Status Definitions (Note 7\)**

* **Power On**: The system is fully initialized and ready to display a source.  
* **Warm Up**: The system is initializing. Do not send other commands during this stage.  
* **Cool Down**: The fan is cooling the lamp after shutdown. Do not send other commands during this stage.  
* **Power Off**: The system is turned off. For rebooting via LAN, the "Standby LAN Control" setting must be set to "ON".

#### **Function-Specific Notes (Notes 8 & 9)**

* **Mute**: The "Mute" function is only active when there is an input source applied.
* **Auto Adjust**: This function is only active with non-digital input sources, such as VGA/Computer/D-sub.
* **Reset to Factory Default**: After sending this command, the user must reboot the projector for the parameters to be cleared.
