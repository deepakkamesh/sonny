/* Package devices manages all the devices on Sonny and provides gRPC interface to interact with it
 */
package devices

type device interface {
	start()
}

type devices struct {
}

// New returns a new initialized devices.
func New() *devices {

}
