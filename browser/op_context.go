package browser

/* OpContext implements context.Context. This is used for tasks such as
 * interrupting operations, and passing in cache invalidation and updation
 * channels.
 */

type OpContext struct {
}
