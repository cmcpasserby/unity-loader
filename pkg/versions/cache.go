package versions

import "time"

type Cache struct {
    Time time.Time
    Versions []ExtendedData
}
