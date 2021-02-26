import {
  Button,
  ButtonGroup,
  Callout,
  Classes,
  Dialog,
  FormGroup,
  HTMLTable,
  InputGroup,
  Tag,
} from "@blueprintjs/core";
import { Popover2 } from "@blueprintjs/popover2";
import React from "react";
import { IModifier } from "types";

let regDec = /^-?(\d+(\.\d+)?|\.\d+)$/;
let regAllowable = /^-?(\d+)?(\.)?(\d+)?$/;

function ModsTable({
  mods,
  handleAdd,
  handleRemove,
}: {
  mods: IModifier[];
  handleAdd: (modIndex: number, value: number, desc: string) => void;
  handleRemove: (modIndex: number, index: number) => void;
}) {
  const [show, setShow] = React.useState<boolean>(false);
  const [modIndex, setModIndex] = React.useState<number>(-1);
  const [label, setLabel] = React.useState<string>("");
  const [value, setValue] = React.useState<string>("0");
  const [desc, setDesc] = React.useState<string>("");

  const handleRemoveWrapper = (modIndex: number, index: number) => {
    return () => handleRemove(modIndex, index);
  };
  const handleAddWrapper = () => {
    if (modIndex >= 0 && modIndex < mods.length) {
      if (regDec.test(value)) {
        handleAdd(modIndex, parseFloat(value), desc);
        handleCloseAdd();
      }
    }
  };
  const handleOpenAdd = (modIndex: number, label: string) => {
    return () => {
      setModIndex(modIndex);
      setShow(true);
      setLabel(label);
    };
  };
  const handleCloseAdd = () => {
    setModIndex(-1);
    setShow(false);
    setLabel("");
  };
  var rows = mods.map((e, i) => {
    var tags = e.list.map((t, j) => {
      return (
        <Popover2 content={t.desc} key={"tag-" + e.type + "-" + j}>
          <Tag interactive onRemove={handleRemoveWrapper(i, j)}>
            {t.value}
          </Tag>
        </Popover2>
      );
    });
    return (
      <tr key={i}>
        <td>
          {e.helper ? (
            <Popover2 content={e.helper}>{e.label}</Popover2>
          ) : (
            e.label
          )}
        </td>
        <td>{tags}</td>
        <td>
          <Button
            small
            minimal
            icon="add"
            onClick={handleOpenAdd(i, e.label)}
          />
        </td>
      </tr>
    );
  });
  return (
    <div>
      <HTMLTable style={{ width: "100%" }} striped>
        <thead>
          <tr>
            <th style={{ minWidth: "200px", width: "10%" }}>Mods</th>
            <th style={{ width: "100%" }}>Percentage</th>
            <th></th>
          </tr>
        </thead>
        <tbody>{rows}</tbody>
      </HTMLTable>
      <Dialog isOpen={show} onClose={handleCloseAdd} title={label}>
        <div className={Classes.DIALOG_BODY}>
          {modIndex < 0 || modIndex >= mods.length ? (
            <Callout intent="warning">Something went wrong</Callout>
          ) : null}

          <FormGroup label="modifier value in decimals">
            <InputGroup
              placeholder="Input value"
              value={value}
              onChange={(e) => {
                if (regAllowable.test(e.target.value)) {
                  setValue(e.target.value);
                }
              }}
            />
          </FormGroup>
          <FormGroup label="description" labelInfo="optional">
            <InputGroup
              placeholder="Description..."
              value={desc}
              onChange={(e) => {
                setDesc(e.target.value);
              }}
            />
          </FormGroup>
          <ButtonGroup fill>
            <Button
              intent="primary"
              onClick={handleAddWrapper}
              disabled={!regDec.test(value)}
            >
              Add
            </Button>
            <Button onClick={handleCloseAdd}>Cancel</Button>
          </ButtonGroup>
        </div>
      </Dialog>
    </div>
  );
}

export default ModsTable;

/**
 * <tr>
          <td>Atk %</td>
          <td>
            <Tag interactive onRemove={() => {}}>
              49.6%
            </Tag>
          </td>
          <td>
            <Button small minimal icon="add" />
          </td>
        </tr>
        <tr>
          <td>Elemental Damage %</td>
          <td>
            <Tag interactive onRemove={() => {}}>
              49.6%
            </Tag>
          </td>
          <td>
            <Button small minimal icon="add" />
          </td>
        </tr>
        <tr>
          <td>Crit Chance %</td>
          <td>
            <Tag interactive onRemove={() => {}}>
              49.6%
            </Tag>
          </td>
          <td>
            <Button small minimal icon="add" />
          </td>
        </tr>
        <tr>
          <td>Crit Damage %</td>
          <td>
            <Tag interactive onRemove={() => {}}>
              49.6%
            </Tag>
          </td>
          <td>
            <Button small minimal icon="add" />
          </td>
        </tr>
        <tr>
          <td>Damage %</td>
          <td>
            <Tag interactive onRemove={() => {}}>
              49.6%
            </Tag>
          </td>
          <td>
            <Button small minimal icon="add" />
          </td>
        </tr>
        <tr>
          <td>EM Increase</td>
          <td>
            <Tag interactive onRemove={() => {}}>
              49.6%
            </Tag>
          </td>
          <td>
            <Button small minimal icon="add" />
          </td>
        </tr>
        <tr>
          <td>Reaction Bonus</td>
          <td>
            <Tag interactive onRemove={() => {}}>
              49.6%
            </Tag>
          </td>
          <td>
            <Button small minimal icon="add" />
          </td>
        </tr>
        <tr>
          <td>Resist Mod</td>
          <td>
            <Tag interactive onRemove={() => {}}>
              49.6%
            </Tag>
          </td>
          <td>
            <Button small minimal icon="add" />
          </td>
        </tr>
        <tr>
          <td>Defence Shred Mod</td>
          <td>
            <Tag interactive onRemove={() => {}}>
              49.6%
            </Tag>
          </td>
          <td>
            <Button small minimal icon="add" />
          </td>
        </tr>
 */
