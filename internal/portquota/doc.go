// Package portquota enforces configurable port binding quotas on a per-subject
// basis, where a subject is typically a process name or OS user.
//
// Usage:
//
//	q := portquota.New(10) // default: 10 ports per subject
//	q.SetLimit("nginx", 50)
//
//	if err := q.Track("nginx", 443); err != nil {
//		log.Println(err)
//	}
//
//	q.Release("nginx", 443)
package portquota
