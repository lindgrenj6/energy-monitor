> Do you have an APSystems ECU?

> Are you sick of the fact that the website doesn't work half the time and the app is....an app?

## WELL HAVE I GOT A SOLUTION FOR YOU

I coudln't tell you what model of ECU I have, all I know is:

1. It was installed on my network when my solar panels were installed
2. Poking at the open ports exposed a web server (❗)
3. It's running PHP and I can get some funky JSON output via a URL

Time to feed this into prometheus so we can get some _actual_ good visualization and history!

The URL is: `http://<ECU IP>/index.php/realtimedata/old_power_graph` which spits out json in this format:
```json
{
  "power": [
    {
      "time": 1663601593000,
      "each_system_power": 2941
    },
    {
      "time": 1663601893000,
      "each_system_power": 4116
    }
  ],
  "today_energy": "19.01",
  "subtitle": ""
}
```

Which translates to roughly this Go schema:
```go
type RealTimeData struct {
	Power []struct {
		Time            int64 `json:"time"`
		EachSystemPower int   `json:"each_system_power"`
	} `json:"power"`
	TodayEnergy string `json:"today_energy"`
	Subtitle    string `json:"subtitle"`
}

```
A few notes:
- `[]Power` is timestamped data in milliseconds since epoch, companies with the power at that time. So it's historical data
- `TodayEnergy` is the amount of power generated today....so it's current data.

#### _Oh boy, a mixture of current AND historical data. My favorite_

Basically what this program does is scrape that data and transform it into a format Prometheus can ingest, something like this:
```prom

# HELP current_power The current power of the system
# TYPE current_power gauge
current_power 3243 1663597393000
current_power 3223 1663597693000
current_power 3355 1663597993000
current_power 4116 1663601893000
current_power 4058 1663602193000
current_power 3405 1663610293000
current_power 3376 1663612693000


# HELP total_power The amount of power generated for the day
# TYPE total_power counter
total_power 19.28
```

The `current_power` metric is exactly what the name suggests - it turns out prometheus understands timestamps if you specify it after the value, neat!

The `total_power` doesn't have a timestamp - so prometheus will pull it and timestamp it for us. Boring but awesome!


### To Run:

1. Build image, push it somewhere (or run it locally)
2. Configure prometheus to scrape wherever it's running whether that is on a box or in a k8s cluster
3. ????
4. Enjoy your green energy a little bit more!
