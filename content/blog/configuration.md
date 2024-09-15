---
title: "Feature Flags"
date: "2024-08-03"
tags: ["Programming"]
---

## Dynamic Configuration

It should always be possible to configure feature flags statically on a per-dimension basis. In addition, dynamic configuration can be used.

- Audit trail.

## Deployment Mechanism

Ideally, all changes should be deployed as part of the CI/CD release process with safety provided by automated tests and monitoring. In reality, this is not always feasible:

- Monitoring gaps result in a bad change being promoted through the pipeline.
- Rollbacks cause significant disruption to the team (along with confusion when trying to recall the bad commit).
- One box environments do not receive sufficient traffic to fully exercise the new code. The smallest deployment unit receives too much traffic to limit blast radius.
- The change is being deployed ahead of its dependencies in which case it cannot be enabled by default.
- The change is being released along with unrelated changes that cannot be easily rolled back.
- The change is too risky for automated deployments and requires manual enablement and validation steps.
- The change is not compatible with all deployment units as a result of a snowflake architecture.

Feature flag based deployments can sometimes strike a better balance between automation and manual action and validation. It will often only be necessary to manually test the feature in a subset of deployment environments after which a static configuration change can be deployed to enable the feature flag by default in all environments.

## Percent Migration

To further reduce the blast radius of a potentially bad change, a percent migration applied to the new code path in a single deployment unit.

This can be implemented in a few different ways for a chosen migration rate, \(r \in [0, 1]\):

- Pick a random number \(x \in [0, 1]\) and check \(x \le p\).
- Maintain a count, \(c\). For each request, check \(c \mod \frac{1}{r} < 1\) and increment \(c\).
- Concatenate the user ID and flag name to obtain \(s\), check \(h(s) \mod 100 < r * 100\).

```go
type RateCounter struct {
	rateProvider func() float64
	counter      atomic.Uint64
}

func (rc *RateCounter) IsOpen() (bool, error) {
	rate := rc.rateProvider()
	if !(0.0 <= rate && rate <= 1.0) {
		return false, fmt.Errorf("Invalid rate %f", rate)
	}

	count := rc.counter.Add(1)
	period := 1.0 / rate

	return math.Mod(float64(count) - 1, period) < 1, nil
}

func main() {
	rc := RateCounter{rateProvider: func() float64 {
		return 0.7
	}}

	total := 0
	for i := 0; i < 100; i++ {
		isOpen, err := rc.IsOpen()
		if err != nil {
			log.Fatal(err)
		}

		if isOpen {
			total++
		}

		log.Println(isOpen)
	}

	log.Println(total)
}
```

```go
// isInPercentage check if the user is in the cohort for the toggle.
func (f *FlagData) isInPercentage(flagName string, user ffcontext.Context) bool {
	percentage := int32(f.getActualPercentage())
	maxPercentage := uint32(100 * percentageMultiplier)

	// <= 0%
	if percentage <= 0 {
		return false
	}
	// >= 100%
	if uint32(percentage) >= maxPercentage {
		return true
	}

	hashID := utils.Hash(flagName+user.GetKey()) % maxPercentage
	return hashID < uint32(percentage)
}
```

<!-- https://github.com/thomaspoignant/go-feature-flag/blob/c84e9326f895c67913f04949c4a76645c18da48f/testutils/flagv1/flag_data.go#L125 -->

## Private Beta

Sometimes we can choose the test cohort explicitly. If this is the case, provide access to the feature as a "private beta" experience for select customers with the understanding that it is not in its final state.

## Scheduled Enablement/Expiration

It should be possible to specify "active after" and "inactive after" dates for feature flags.


## Feature Flag Hygiene

Feature flags should be regularly reviewed and removed from code where possible.

Feature flags require visibility to be useful; otherwise, they become a liability.


## Toothless Mode

There exists a trade-off between the extent to which the new code path runs and how toothless that new code path actually is.


## De-fanged Side-effects


## Common Pitfalls

- Don't read the value of a feature flag twice on the same code path (unless this is explicitly intended). If the flag is set dynamically, then it's value may have changed between reads with unpredictable results.


