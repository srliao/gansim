**TL;DR:** It would take on avg 215 5\* artifact drops (std dev 120.92) or roughly 24 days non stop domain farming to gear a F2P 4pc BS Ganyu to 99th percentile damage output. While this is specific to Ganyu, I think the results should be fairly extensible to other characters as just about any main DPS will want to run Atk %/Ele %/CC or CD.

## Original motivation

I originally wrote this sim as a way to compare to damage output of Amos R1 vs Prototype R5 just to see how big of a gap there is. There were plenty of spreadsheet calculations but those all had assumed certain artifact stats. Obviously not everyone will have the same stats as those assumptions and I wasn't sure if having different stats will significantly affect the outcome.

Instead, I turned to what I thought would make more sense. A monte carlo simulation of randomily drawn artifacts. These randomily drawn artifacts would have the correct main stat (i.e. in my case, Atk % Sands, Ele % Goblet, and Crit Dmg Circlet) and randomily drawn substats. This would essentially simulate a situation where a player gets the correct artifact with the correct main stat from the artifact domain, and then upgrades that to +20, disregarding whatever substat is on said artifact.

By doing this doing this enough times (it converges pretty well even at 100,000 simulations), we can then plot out a histogram showing the distribution of damage given correct main stat and random substat (and upgrades).

Here's an example of the output

https://imgur.com/p9INyM3

And then after some discussion with some friends, I thought wouldn't it be great if I can then figure out on average how long would it take to farm up the artifacts required to achieve top dps.

So using the result of the sim above (which tells me for example that 95% of the possible damage output is below **x**), I create another sim that emulates drawing random artifacts and figuring out how many draws it took to achieve a certain damage threshold.

## Algorithm

The steps to generate a random +20 artifact with a given main stat is roughly as follows:

```
    roll chance of getting 3 vs 4 lines initially
    generate 4 random substats (see assumptions for rates)
    let number of upgrades = 5
    if initially = 3 lines
        number of upgrades = 4 (since +4 gives us the 4th line)
    roll upgrades
```

The steps to sim out the damage distribution is roughly as follows:

```
Repeat n times:
1. generate a set of artifacts with atk % sands, ele % goblet, and crit dmg circlet main stat and random substat, upgraded to +20
2. calculate average damage
```

And the follow to simulate farming for artifacts

```
Repeat while max damage < threshold:
    roll random slot
    roll 50/50 for on-set vs off-set
    if !on-set && slot != goblet
        discard, continue with next draw
    roll random main stat
    if main stat = ele %, roll 1/6 for correct element
        if not correct element
            discard, continue with next draw
    if !on-set && main stat != ele%
        discard, continue with next draw
    if slot == sands, goblet, or circlet
        if main stat != atk % or ele % or crit % or crit dmg
            discard, continue with next draw
    //at this point we can only have either offset goblet with correct +ele %, or atk % sands, or atk %/cc/cd circlet
    //note that the weights are set such that you can't roll a sands with cc/cd
    roll random substats, upgrade to +20 (because sometimes you get 3 liners so you need to take it to +4 at lest)
    //here I consider good substat to be cc/cd/atk%/atk/er, must = atk%/cc/cd
    if # of good substat < 2 and number of must substat < 1
        discard, continue with next draw
    //at this point this is a keeper
    caculate damage with this new artifact with existing other slots
    if damage > existing
        max = damage
        replace existing slot with this new one
```

## Probability Weights & Assumptions

For probabilities, I used data from the data gathering discord (without which this would not have been possible). There are some limitation with the data set. Unfortunately the sample size is relatively small for rarer items such as cc/cd circlets. I'm not stats guy so I have no idea how much error this will introduce to the result but I imagine this current result should still provide a good estimate.

Some other assumptions had to be made as well:

- I assume there's equal chance of upgrading any of the 4 substats
- There are 4 tiers of possible substat upgrades. I assume it's equal chance of getting any of the 4 tiers on an upgrade
- I assume that the probability of getting a substat is dependent on the main stat (a forced assumption because otherwise I have no idea how to find the probability of getting for example atk % on a crit circlet)
- For substat probability, I assume that if let's say there are 50 5\* Atk % Sands in the sample. Out of these 50 artifacts, we have a total of 160 substats. Out of the 160 substats, 12 rolled crit rate, 15 rolled crit damage, 20 rolled flat attack. Then the probability of getting crit rate on a atk % sand is 12/160. If the first substat is crit rate, then the probability that the second substat is crit damage would be 15/(160-12). Here I adjust the denominator as crit rate is no longer in the pool.

Would appreciate anyone that can help validate these assumptions and/or come up with better ways to arrive at the probabilities/simulate artifact generation

For damage formula, I used the formulate posted [here](https://old.reddit.com/r/Genshin_Impact/comments/krg2ic/the_complete_genshin_impact_damage_formula/). However I don't use the full formula as I'm not considering reactions.

## Results

Here's the damage simulation for a Lvl 90 Ganyu with Lvl 90 Prototype R1 (procced), 4pc BS vs a lvl 88 Lawachurl. Note that this is the combined average damage of the arrow + bloom on charged level 2. Note that this is with C1 Ganyu (which increases the damage of the bloom re 15% resist reduction) I'm using these specific values because that's what I tested against in game to make sure my damage calculation is correct.

https://imgur.com/3TagG16

I only did 1mil simulation since it converges pretty nicely already. 10mil takes too long.

From this simulation we can calculate the 95th percentile is roughly 20,200 damage. So using that, I sim how many artifact drops it would take to reach that damage threshold

https://imgur.com/pqTiTws

Again 1mil simulation because my code is slow and 10mil would take forever. Here's the summary statistics:

Min: 7
Max: 1400
Mean: 204.79
Median: 170 (an estimate becaused I binned the data first increments of 10 before calculating median)
Std Dev: 122.94

## Other remarks

Looking at these results, I guess it's really not all that bad to gear a character. Considering that a good set of +20 artifacts is actually transferable between characters. In addition, if it actually takes on average 200 runs to gear 1 character then I can probably start slacking off on my daily artifact runs. I believe it takes roughly 18 runs to level one 5\* to +20 so I should have enough exp just from farming the domain itself.

Also, I'm quite happy with how the sim turned out. I think the first damage sim is a much more useful way of comparing weapons vs. the traditional spreadsheet. This is because I make no assumption about the substats as that could play a big role in the overall damage. Much easier to look at and compare two distributions to determine which weapon is better.

The sim can also be extended to other character very easily as well by simply modifying the profile. If I have time I probably want to try out some other characters I own.

Finally, I'm not really a math/stat/cs guy so would appreciate anyone that can help look over this for me for any inaccuracies.
