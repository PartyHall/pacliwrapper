# Pulseaudio-cli Wrapper

## Why

I didn't want to use CGO to access Pipewire / Pulse Audio.

I did try to access pulseaudio by socket and it was a pain, it never answered me, nothing really worked.

This is ugly, I don't like it, but it does the job.

It is intentionally very basic, the only purpose is for PartyHall to get & monitor the volume, set it, and change the output device. No other feature will be implemented.

## Requirements
- A pulseaudio server running
- The pactl command

## Usage

```go
pa := pacliwrapper.New()
if err := pa.Refresh(); err != nil {
    fmt.Println("Failed to refresh the list of sinks!")
    fmt.Println(err)
    os.Exit(1)
}

if pa.MainDevice == nil {
    fmt.Println("Failed to find main device!")
    os.Exit(1)
}

fmt.Println(pa.MainDevice.GetVolume())
fmt.Println(pa.MainDevice.SetVolume(30)) // Set the volume to 30%

for _, d := range pa.Devices {
    if d.Name == "alsa_output.pci-0000_0b_00.4.analog-stereo" {
        pa.SetDefaultOutput(d)
        break
    }
}
```

## License
Copyright © PartyHall team

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
The Software is provided “as is”, without warranty of any kind, express or implied, including but not limited to the warranties of merchantability, fitness for a particular purpose and noninfringement. In no event shall the authors or copyright holders be liable for any claim, damages or other liability, whether in an action of contract, tort or otherwise, arising from, out of or in connection with the software or the use or other dealings in the Software.
