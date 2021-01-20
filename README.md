# gansim

## Usage

Make sure you have [go](https://golang.org/doc/install) and [git](https://git-scm.com/) installed.

Then

```
git clone https://github.com/srliao/gansim.git
cd ./gansim
mkdir graphs
go run main.go
```

By default the graphs will be generated `./graphs` so make sure you create a graphs folder or else this will fail.

## Config explanation

There are two config files, the main `config.yaml` and the various profile.

`config.yaml` currently have the following parameters:

| Param             | Explanation                                                                                                                                                                                            |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `GraphOutput`     | Tells the script where to save the graph output                                                                                                                                                        |
| `NumWorker`       | Tells the script how many goroutines (not quite threads) to spawn. I have a AMD 3600X and I find the best number for me is 24                                                                          |
| `NumSim`          | Number of simulation to run. Starts converging by 100k and converges pretty well by 1mil                                                                                                               |
| `BinSize`         | The size of the bin for the histograph/graph output. Stick to 100 if not sure                                                                                                                          |
| `WriteToCSV`      | Dump out the histogram info to CSV file. CSV file name specified in each profile config separately                                                                                                     |
| `DamageType`      | Type of damage to report on. Accepted values are `normal` `avg` and `crit`. If this is not specified or doesn't match the spelling/case exactly it'll assume `avg`                                     |
| `MainStatFile`    | Path to the main stat by lvl (csv file). just leave this as default                                                                                                                                    |
| `SubstatTierFile` | This is the amount of stat you gain per substat lvl up by tier (csv file). Just leave as default                                                                                                       |
| `Profiles`        | Here you specify the profiles to run for the sim. It's an array so you can run as many at once as you like but keep in mind it'll just take longer and you may end up with lots of lines on your graph |

Each profile have the following fields:

| Param              | Explanation                                                                                                                                    |
| ------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| `Output`           | tells the script what to name the csv file if we're writing to csv                                                                             |
| `Label`            | label for this profile, shows up on the graph                                                                                                  |
| `CharLevel`        | Character level                                                                                                                                |
| `CharBaseAtk`      | Character base attack. Script does not calculate for you so you need to use the right amount of the right level                                |
| `WeaponBaseAtk`    | Weapon base attack. Make sure to use the right base attack for whatever level the weapon is                                                    |
| `EnemyLevel`       | Level of enemy to sim against. Affects their resistance                                                                                        |
| `ArtifactMaxLevel` | What level to upgrade the artifacts to. Useful to compare output at +16 vs +20 for example. +4 and lower are not tested. Use at your own risk  |
| `Sands`            | Main stat type of the sand. Accepted values are `DEF% DEF HP HP% ATK ATK% ER EM CR CD Heal Ele% Phys%`. Script does not perform validity check |
| `Goblet`           | Same as `Sands`                                                                                                                                |
| `Circlet`          | Same as `Sands`                                                                                                                                |
| `SubstatFile`      | This file specify the substat weightings. Take a look at the file on this repo as an example. Must be in this format (including the header)    |
| `Abilities`        | This is an array specifying which abilities to calculate damage for                                                                            |

Each item in `Abilities` have the following fields:

| Param       | Explanation                                                                                                                                                                                          |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Talent`    | This is the talent damage % in decimals                                                                                                                                                              |
| `AtkMod`    | List of attack % modifiers in decimals. For example, weapon substat, any weapon proc effects, etc.. should all be on here                                                                            |
| `EleMod`    | Same as `AtkMod`. This should include any Ele% modes that's not part of the artifact main/substat, **including** any set bonuses                                                                     |
| `CCMod`     | List of critical rate increase/decrease mods. Make sure to include anything that's not part of artifact main/substat here **including** character base crit rate, set bonus, resonance bonus, etc... |
| `CDMod`     | see `CCMod` but for crit damage                                                                                                                                                                      |
| `DmgMod`    | see `CCMod` but for % damage increase (such as Amos bow passive)                                                                                                                                     |
| `ResistMod` | see `CCMod` but for resistance modifier such as superconduct                                                                                                                                         |
| `DefShred`  | see `CCMod` but for defence shredding abilities                                                                                                                                                      |

## Future todo

- refactor out dmg calculation into separate package and write unit tests
- split up character profile and abilities profile so it's easy to test same character but diff rotations
- add reactions

## About

I originally wrote this sim as a way to compare to damage output of Amos R1 vs Prototype R5 just to see how big of a gap there is. There were plenty of spreadsheet calculations but those all had assumed certain artifact stats. Obviously not everyone will have the same stats as those assumptions and I wasn't sure if having different stats will significantly affect the outcome.

Instead, I turned to what I thought would make more sense. A monte carlo simulation of randomily drawn artifacts. These randomily drawn artifacts would have the correct main stat (i.e. in my case, Atk % Sands, Ele % Goblet, and Crit Dmg Circlet) and randomily drawn substats. This would essentially simulate a situation where a player gets the correct artifact with the correct main stat from the artifact domain, and then upgrades that to +20, disregarding whatever substat is on said artifact.

By doing this doing this enough times (it converges pretty well even at 100,000 simulations), we can then plot out a histogram showing the distribution of damage given correct main stat and random substat (and upgrades).

To me this was way more useful than damage calculations on a spreadsheet because not only does this show the min/max/avg scenario, it also shows the variance. It makes more sense to compare two scenarios with the min/max/avg/variance than just a simple cherry picked stat.

Here is an example output showing the distribution of damage for Amos R1 vs various Prototype bos:

![bow comparison](/example.png?raw=true)

## Assumptions

In order to perform a monte carlo simulation, certain assumptions had to be made (whether true or not). In particularly, there has to be assumption as to the probability of rolling each substat and the probability of upgrading a substat.

### Generating artifacts

I generate random artifact line by line. Meaning first I randomly rolled to figure out which substat are present on a +0 artifact (i.e is it 3 lines or 4 lines, probability based on actual samples collected by data gathering discord, and which substats).

Then the +0 artifact is upgraded to the specified level. Every 4 level a random substat is picked and a random upgrade amount is picked.

### Rolling substat

For my simulations, I used the samples collected by the Data Gathering discord. In order to calculate the probability of getting any one particular substat, I would filter the samples collect down to the 5\* artifact of the same slot and same main stat.

For example (these are made up numbers), let's say there are 50 5\* Atk % Sands in the sample. Out of these 50 artifacts, we have a total of 160 substats. Out of the 160 substats, 12 rolled crit rate, 15 rolled crit damage, 20 rolled flat attack.

So then, for my simulation, the probability that the first substat would be crit rate is 12/160. If the first substat rolled into crit rate, then the probability that the second substat is crit damage would be 15/(160-12). Here I adjust the denominator as crit rate is no longer in the pool.

This kind of assumption is basically assuming that the chance of rolling crit rate is independent of crit damage. There is also the question that if 12/160 is even a good way of estimating the probability of crit rate in the first place. Unfortunately I don't have the anwer to that question. I'm not very good at statistics. In fact, most of this I had to ask a stat friend for verification - and the answer I got was: it's probably not exactly right but good enough for now. So until someone can provide me with a better way of estimating the probability, this is what I ended up using.

### Which substat to upgrade

Here I assume there's an equal chance of upgrading any of the 4 substats. This may not be correct as substat upgrades may be weighted as well but we do not have any data to suggest one way or another.

### How much to upgrade the substat by

There are for possible tiers of upgrades. I assume that each of the 4 tiers have equal chance of happening. Again, this realistically may be weighted as well but we do not have any data to suggest one way or another.

### Damage formula

I use the damage formula posted on Reddit [here](https://old.reddit.com/r/Genshin_Impact/comments/krg2ic/the_complete_genshin_impact_damage_formula/). However, I don't use the full formula currently as I don't include any of the reactions.

## Possible other use case

Even though this was originally written to compare specifically the damage output of Amos R1 vs Prototype for Ganyu, the simulation can be applied to other situations pretty simply by adjusting the config file. See config file explanation for things you can change.

In addition, since atk%/cc/cd/etc... mods can be applied to each individual abilitiy, you can actually sim a rotation of abilities, especially one's that apply special effects such as increase crit chance when affected by cryo. In the first ability you can specify a lower crit rate to simulate that the first ability applies the cryo.
