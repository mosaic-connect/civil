// Package civil provides types for representing civil dates, and times.
//
// Sometimes it is useful to be able to represent a date or time without reference
// to an instance in time within a timezone.
//
// For example, when recording a person's date of birth all that is needed is a date.
// There is no requirement to specify an instant in time within a timezone.
//
// There are also circumstances where an event will be scheduled for a date and time
// in the local timezone, whatever that may be. And example of this might be a schedule for
// taking medication.
package civil
