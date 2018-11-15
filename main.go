package fedrampup
import (
  "log"
  "encoding/csv"
)
func main() {
  config = NewConfig()
  instances, err := NewFetcher(config).Run()
  if err != nil {
    log.Fatal(err)
  }

  w := csv.NewWriter(os.Stdout)
  w.write(Headers)
  for instance := range instances {
    if err := w.Write(instance.Row()); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
  }
  w.Flush()
}
