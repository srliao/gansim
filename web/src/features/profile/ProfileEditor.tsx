import React from "react";
import produce from "immer";
import {
  Button,
  ButtonGroup,
  Card,
  Collapse,
  EditableText,
  FormGroup,
  H3,
  H5,
  H6,
  HTMLSelect,
  HTMLTable,
  InputGroup,
  ITreeNode,
  Switch,
  Tag,
  Tree,
} from "@blueprintjs/core";
import ModsTable from "./ModsTable";
import {
  IAbility,
  IProfile,
  STAT_TYPE_ATKP,
  STAT_TYPE_CD,
  STAT_TYPE_CR,
  STAT_TYPE_DEFP,
  STAT_TYPE_ELEP,
  STAT_TYPE_EM,
  STAT_TYPE_ER,
  STAT_TYPE_HEAL,
  STAT_TYPE_HPP,
  STAT_TYPE_PHYP,
} from "types";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "app/store";
import { pEdit } from "./profileSlice";
import { defaultMods } from "./defaults";

let validNum = /^\d+$/;
let regDec = /^-?(\d+(\.\d+)?|\.\d+)$/;
let regAllowable = /^-?(\d+)?(\.)?(\d+)?$/;

//HP%, DEF%, ATK%, +EM, ER
const sandOpt = [
  { label: "HP%", value: STAT_TYPE_HPP },
  { label: "DEF%", value: STAT_TYPE_DEFP },
  { label: "ATK%", value: STAT_TYPE_ATKP },
  { label: "EM", value: STAT_TYPE_EM },
  { label: "ER", value: STAT_TYPE_ER },
];
//HP%, DEF%, ATK%, +EM, ELE%, Phy%
const gobletOpt = [
  { label: "HP%", value: STAT_TYPE_HPP },
  { label: "DEF%", value: STAT_TYPE_DEFP },
  { label: "ATK%", value: STAT_TYPE_ATKP },
  { label: "EM", value: STAT_TYPE_EM },
  { label: "Ele%", value: STAT_TYPE_ELEP },
  { label: "Phy%", value: STAT_TYPE_PHYP },
];
//HP%, DEF%, ATK%, +EM, CR, CD, Heal
const circletOpt = [
  { label: "HP%", value: STAT_TYPE_HPP },
  { label: "DEF%", value: STAT_TYPE_DEFP },
  { label: "ATK%", value: STAT_TYPE_ATKP },
  { label: "EM", value: STAT_TYPE_EM },
  { label: "CR", value: STAT_TYPE_CR },
  { label: "CD", value: STAT_TYPE_CD },
  { label: "Healing Bonus", value: STAT_TYPE_HEAL },
];

type IAction =
  | { t: "label"; v: string }
  | { t: "character/level"; v: number }
  | { t: "character/base_atk"; v: number }
  | { t: "character/mod/add"; i: number; v: number; d: string }
  | { t: "character/mod/rm"; i: number; j: number }
  | { t: "weapon/base_atk"; v: number }
  | { t: "enemy/level"; v: number }
  | { t: "enemy/phy_resist"; v: number }
  | { t: "enemy/ele_resist"; v: number }
  | { t: "artifact_level"; v: number }
  | { t: "artifact/sand_main"; v: string }
  | { t: "artifact/goblet_main"; v: string }
  | { t: "artifact/circlet_main"; v: string }
  | { t: "abilities/add" }
  | { t: "abilities/multiplier"; a: number; v: number }
  | { t: "abilities/is_vape_melt"; a: number; v: boolean }
  | { t: "abilities/is_physical"; a: number; v: boolean }
  | { t: "abilities/vape_melt_multiplier"; a: number; v: number }
  | { t: "abilities/mod/add"; a: number; i: number; v: number; d: string }
  | { t: "abilities/mod/rm"; a: number; i: number; j: number }
  | { t: "abilities/rm"; a: number };

function ProfileEditor({ index }: { index: number }) {
  const { profiles } = useSelector((state: RootState) => {
    return {
      profiles: state.profile.profiles,
    };
  });
  const dispatch = useDispatch();
  const [showCharMod, setShowCharMod] = React.useState<boolean>(false);

  if (index < 0 || index >= profiles.length) {
    return <div>Unexpected error; corrupt profile. Select a different one</div>;
  }

  const p: IProfile = profiles[index];

  const handleEdit = (a: IAction) => {
    dispatch(
      pEdit({
        index: index,
        p: produce(p, (next) => {
          switch (a.t) {
            case "label":
              next.label = a.v;
              return;
            case "character/level":
              next.character.level = a.v;
              return;
            case "character/base_atk":
              next.character.base_atk = a.v;
              return;
            case "character/mod/add":
              if (a.i >= 0 && a.i < next.character.mods.length) {
                next.character.mods[a.i].list.push({ value: a.v, desc: a.d });
              }
              return;
            case "character/mod/rm":
              if (a.i >= 0 && a.i < next.character.mods.length) {
                if (a.j >= 0 && a.j < next.character.mods[a.i].list.length) {
                  next.character.mods[a.i].list.splice(a.j, 1);
                }
              }
              return;
            case "weapon/base_atk":
              next.weapon.base_atk = a.v;
              return;
            case "enemy/level":
              next.enemy.level = a.v;
              return;
            case "enemy/ele_resist":
              next.enemy.ele_resist = a.v;
              return;
            case "enemy/phy_resist":
              next.enemy.phy_resist = a.v;
              return;
            case "artifact/sand_main":
              next.artifact_main_stats.SANDS = a.v;
              return;
            case "artifact/goblet_main":
              next.artifact_main_stats.GOBLET = a.v;
              return;
            case "artifact/circlet_main":
              next.artifact_main_stats.CIRCLET = a.v;
              return;
            case "abilities/add":
              next.abilities.push({
                id: Date.now(),
                multiplier: 0,
                is_vape_melt: false,
                vape_melt_multiplier: 1.5,
                mods: defaultMods(),
                is_physical: false,
              });
              return;
            case "abilities/multiplier":
              if (a.a < 0 || a.a >= next.abilities.length) return;
              next.abilities[a.a].multiplier = a.v;
              return;
            case "abilities/is_physical":
              if (a.a < 0 || a.a >= next.abilities.length) return;
              next.abilities[a.a].is_physical = a.v;
              return;
            case "abilities/is_vape_melt":
              if (a.a < 0 || a.a >= next.abilities.length) return;
              next.abilities[a.a].is_vape_melt = a.v;
              return;
            case "abilities/vape_melt_multiplier":
              if (a.a < 0 || a.a >= next.abilities.length) return;
              next.abilities[a.a].vape_melt_multiplier = a.v;
              return;
            case "abilities/mod/add":
              if (a.a < 0 || a.a >= next.abilities.length) return;
              if (a.i < 0 || a.i >= next.abilities[a.a].mods.length) return;
              next.abilities[a.a].mods[a.i].list.push({
                value: a.v,
                desc: a.d,
              });
              return;
            case "abilities/mod/rm":
              if (a.a < 0 || a.a >= next.abilities.length) return;
              if (a.i < 0 || a.i >= next.abilities[a.a].mods.length) return;
              if (a.j < 0 && a.j >= next.abilities[a.a].mods[a.i].list.length)
                return;
              next.abilities[a.a].mods[a.i].list.splice(a.j, 1);
              return;
            case "abilities/rm":
              if (a.a < 0 || a.a >= next.abilities.length) return;
              next.abilities.splice(a.a, 1);
              return;
            default:
              return;
          }
        }),
      })
    );
  };

  let abilities = p.abilities.map((a, i) => {
    return (
      <AbilityCard key={a.id} index={i} ability={a} handleEdit={handleEdit} />
    );
  });

  return (
    <div>
      <H3>
        <EditableText
          value={p.label}
          onChange={(v) =>
            handleEdit({
              t: "label",
              v: v,
            })
          }
        ></EditableText>
      </H3>
      <Card style={{ marginTop: "10px" }}>
        <H5>Basic Info</H5>
        <div className="row">
          <div className="col-xs">
            <FormGroup
              helperText="Level of the character"
              label="Character Level"
            >
              <InputGroup
                placeholder="Enter character level"
                value={p.character.level.toString()}
                onChange={(e) => {
                  if (!validNum.test(e.target.value)) return;
                  handleEdit({
                    t: "character/level",
                    v: parseInt(e.target.value),
                  });
                }}
              />
            </FormGroup>
          </div>
          <div className="col-xs">
            <FormGroup
              helperText="Character base attack"
              label="Character Base Attack"
            >
              <InputGroup
                placeholder="Enter character base attack"
                value={p.character.base_atk.toString()}
                onChange={(e) => {
                  if (!validNum.test(e.target.value)) return;
                  handleEdit({
                    t: "character/base_atk",
                    v: parseInt(e.target.value),
                  });
                }}
              />
            </FormGroup>
          </div>
          <div className="col-xs">
            <FormGroup
              helperText="Weapon base at tack"
              label="Weapon Base Attack"
            >
              <InputGroup
                placeholder="Enter weapon base attack "
                value={p.weapon.base_atk.toString()}
                onChange={(e) => {
                  if (!validNum.test(e.target.value)) return;
                  handleEdit({
                    t: "weapon/base_atk",
                    v: parseInt(e.target.value),
                  });
                }}
              />
            </FormGroup>
          </div>
        </div>
        <div className="row">
          <div className="col-xs">
            <FormGroup helperText="Level of enemy" label="Enemy Level">
              <InputGroup
                placeholder="Enter enemy level"
                value={p.enemy.level.toString()}
                onChange={(e) => {
                  if (!validNum.test(e.target.value)) return;
                  handleEdit({
                    t: "enemy/level",
                    v: parseInt(e.target.value),
                  });
                }}
              />
            </FormGroup>
          </div>
          <div className="col-xs">
            <FormGroup helperText="Enemy ele resist" label="Ele resist">
              <InputGroup
                placeholder="Enter elemental resist"
                value={p.enemy.ele_resist.toString()}
                onChange={(e) => {
                  if (!regAllowable.test(e.target.value)) return;
                  handleEdit({
                    t: "enemy/ele_resist",
                    v: parseFloat(e.target.value),
                  });
                }}
              />
            </FormGroup>
          </div>
          <div className="col-xs">
            <FormGroup helperText="Enemy phys resist" label="Phy resist">
              <InputGroup
                placeholder="Enter physical resist"
                value={p.enemy.phy_resist.toString()}
                onChange={(e) => {
                  if (!regAllowable.test(e.target.value)) return;
                  handleEdit({
                    t: "enemy/phy_resist",
                    v: parseFloat(e.target.value),
                  });
                }}
              />
            </FormGroup>
          </div>
        </div>
        <Button onClick={() => setShowCharMod(!showCharMod)}>
          {showCharMod
            ? "Hide character modifiers"
            : "Show character modifiers"}
        </Button>
        <Collapse isOpen={showCharMod}>
          Enter in the table below any stats modifications that will apply to
          all abilities. Use this for character ascension stats, weapon
          substats, enemy resistance, etc..
          <ModsTable
            mods={p.character.mods}
            handleAdd={(modIndex, val, desc) =>
              handleEdit({
                t: "character/mod/add",
                i: modIndex,
                v: val,
                d: desc,
              })
            }
            handleRemove={(modIndex, index) =>
              handleEdit({
                t: "character/mod/rm",
                i: modIndex,
                j: index,
              })
            }
          />
        </Collapse>
      </Card>
      <Card style={{ marginTop: "10px" }}>
        <H5>Artifact Main Stats</H5>
        <div className="row">
          <div className="col-xs">
            <FormGroup helperText="Helper text with details..." label="Sand">
              <HTMLSelect
                options={sandOpt}
                value={p.artifact_main_stats.SANDS}
                fill
                onChange={(e) => {
                  handleEdit({
                    t: "artifact/sand_main",
                    v: e.currentTarget.value,
                  });
                }}
              />
            </FormGroup>
          </div>
          <div className="col-xs">
            <FormGroup helperText="Helper text with details..." label="Goblet">
              <HTMLSelect
                options={gobletOpt}
                value={p.artifact_main_stats.GOBLET}
                fill
                onChange={(e) => {
                  handleEdit({
                    t: "artifact/goblet_main",
                    v: e.currentTarget.value,
                  });
                }}
              />
            </FormGroup>
          </div>
          <div className="col-xs">
            <FormGroup helperText="Helper text with details..." label="Circlet">
              <HTMLSelect
                options={circletOpt}
                value={p.artifact_main_stats.CIRCLET}
                fill
                onChange={(e) => {
                  handleEdit({
                    t: "artifact/circlet_main",
                    v: e.currentTarget.value,
                  });
                }}
              />
            </FormGroup>
          </div>
        </div>
      </Card>
      <Card style={{ marginTop: "10px" }}>
        <H5>Abilities</H5>
        {abilities.length === 0 ? "No abilities. Add one first" : abilities}
        <ButtonGroup fill>
          <Button
            intent="primary"
            onClick={() => handleEdit({ t: "abilities/add" })}
          >
            Add Ability
          </Button>
        </ButtonGroup>
      </Card>
    </div>
  );
}

function AbilityCard({
  ability,
  index,
  handleEdit,
}: {
  ability: IAbility;
  index: number;
  handleEdit: (a: IAction) => void;
}) {
  const [showCharMod, setShowCharMod] = React.useState<boolean>(false);
  const [mul, setMul] = React.useState<string>(ability.multiplier.toString());
  return (
    <Card elevation={3} style={{ margin: "20px" }}>
      <H6>
        <EditableText value={"ability name"}></EditableText>
      </H6>
      <div className="row">
        <div className="col-xs">
          <FormGroup
            helperText="Talent damage % of the ability in decimals"
            label="Ability/Talent Damage %"
          >
            <InputGroup
              placeholder="Enter damage % in decimal format"
              value={mul}
              intent={isNaN(ability.multiplier) ? "danger" : "none"}
              onChange={(e) => {
                if (!regAllowable.test(e.target.value)) return;
                setMul(e.target.value);
                handleEdit({
                  t: "abilities/multiplier",
                  a: index,
                  v: parseFloat(e.target.value),
                });
              }}
            />
          </FormGroup>
        </div>
        <div className="col-xs">
          <FormGroup
            helperText="Damage type of the ability. Typically physical for normal attack or elemental for skill/burst"
            label="Damage Type"
          >
            <Switch
              innerLabel={"Elemental"}
              innerLabelChecked={"Physical"}
              onChange={(e) => {
                handleEdit({
                  t: "abilities/is_physical",
                  a: index,
                  v: e.currentTarget.checked,
                });
              }}
            />
          </FormGroup>
        </div>
        <div className="col-xs">
          <FormGroup
            helperText="Toggle to apply melt/vaporize bonus to this ability"
            label="Melt/Vaporize"
          >
            <Switch
              innerLabel={"No Melt/Vaporize"}
              innerLabelChecked={"Melt/Vaporize Applied"}
              onChange={(e) => {
                handleEdit({
                  t: "abilities/is_vape_melt",
                  a: index,
                  v: e.currentTarget.checked,
                });
              }}
            />
          </FormGroup>
        </div>
        <div className="col-xs">
          <FormGroup
            helperText="Modifier for melt/vaporize. 2x if Fire on Cryo or Hydro on Fire. Other 1x"
            label="Vaporize Modifier"
          >
            <Switch
              innerLabel={"1x"}
              innerLabelChecked={"2x"}
              onChange={(e) => {
                handleEdit({
                  t: "abilities/vape_melt_multiplier",
                  a: index,
                  v: e.currentTarget.checked ? 2 : 1,
                });
              }}
            />
          </FormGroup>
        </div>
      </div>
      <Button onClick={() => setShowCharMod(!showCharMod)}>
        {showCharMod ? "Hide ability modifiers" : "Show ability modifiers"}
      </Button>
      <Collapse isOpen={showCharMod}>
        Enter in table below any stat modification that <b>only</b> applies to
        the current ability
        <ModsTable
          mods={ability.mods}
          handleAdd={(modIndex, val, desc) =>
            handleEdit({
              t: "abilities/mod/add",
              a: index,
              i: modIndex,
              v: val,
              d: desc,
            })
          }
          handleRemove={(modIndex, index) =>
            handleEdit({
              t: "abilities/mod/rm",
              a: index,
              i: modIndex,
              j: index,
            })
          }
        />
      </Collapse>
      <br />
      <Button
        fill
        onClick={() =>
          handleEdit({
            t: "abilities/rm",
            a: index,
          })
        }
        icon="delete"
        intent="danger"
      >
        Delete This Ability
      </Button>
    </Card>
  );
}

export default ProfileEditor;
