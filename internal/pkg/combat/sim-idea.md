# Idea for full combat sim

```
type Enemy struct {
    //enemy basic info
    Level
    Res

    HP //how much hp left on this target
    ActiveAura []Aura //aura applied to this enemy
    StatusEffects []Effects //any other effects maybe? like Omen debuff?
}

type Aura struct {
    //whatever aura
    //ticks remaining in fps
    //strength
}

type Effect struct {
    //some kind of structure to describe how to handle effect calc
    //could just be straight up + dmg or could be something else
    //probably use some sort of constant for the engine to parse
}

type Ability struct {
    Element
    CD //in frames how long this should be
    Execution //in frames
    AvailableIn //how many frames left before this ability is available -> should be execution + CD on first execution, -1 per tick

}

type Character struct {
    whatever you need for base stat goes here
}

type Artifact struct {
    whatever to describe the artifacts
}

```

- runs on 60fps. sim checks what to do every frame
- sim would loop while HP > 0 and frames remain > 0
- action list should be a combination of specific rotations and a priority list
- rotation would be something like Mona burst where you want to apply certain skills in a very specific order
- a rotation is just a combination skill in a priority list

???

- how to deal with tick damage like electro charged? or venti ult? maybe add a DOT effect?

How to structure priority list?
