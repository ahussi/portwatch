// Package resolver provides port-to-service name resolution for portwatch.
//
// It ships with a built-in table of well-known IANA service names and allows
// callers to register custom mappings at runtime. The Resolver is safe for
// concurrent use.
//
// Usage:
//
//	r := resolver.New()
//	fmt.Println(r.Name(443))   // "https"
//	fmt.Println(r.Name(9999))  // "port/9999"
//
//	r.Register(9090, "prometheus", "tcp")
//	fmt.Println(r.Name(9090))  // "prometheus"
package resolver
