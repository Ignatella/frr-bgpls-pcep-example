In order to use pcep frr should be patched:

File: `pcep_msg_obcjets.h`

```c
enum pcep_ro_subobj_types {
	RO_SUBOBJ_TYPE_IPV4 = 1,  /* RFC 3209 */
	RO_SUBOBJ_TYPE_IPV6 = 2,  /* RFC 3209 */
	RO_SUBOBJ_TYPE_LABEL = 3, /* RFC 3209 */
	RO_SUBOBJ_TYPE_UNNUM = 4, /* RFC 3477 */
	RO_SUBOBJ_TYPE_ASN = 32,  /* RFC 3209, Section 4.3.3.4 */
	// PATCHED 36 -> 5
	RO_SUBOBJ_TYPE_SR = 5, /* RFC 8408, draft-ietf-pce-segment-routing-16.
				   Type 5 for draft07 has been assigned to
				   something else. */
	RO_SUBOBJ_UNKNOWN
};
```
