import React from "react";
import {
  Button,
  ButtonGroup,
  ControlGroup,
  Divider,
  H3,
  H4,
  H5,
  InputGroup,
  ITreeNode,
  Tree,
  TreeNode,
} from "@blueprintjs/core";
import ModsTable from "./ModsTable";
import ProfileEditor from "./ProfileEditor";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "app/store";
import { IProfile } from "types";
import { load, pAdd, save } from "./profileSlice";

function Profile() {
  const { profiles, isSaved } = useSelector((state: RootState) => {
    return {
      profiles: state.profile.profiles,
      isSaved: state.profile.saved,
    };
  });
  const dispatch = useDispatch();
  const [selected, setSelect] = React.useState<number>(-1);
  const [newProfile, setNewProfile] = React.useState<string>("");

  const handleNewProfile = () => {
    dispatch(pAdd({ label: newProfile }));
  };

  const handleNodeClick = (nodeData: ITreeNode) => {
    if (typeof nodeData.id === "number") setSelect(nodeData.id);
  };

  var treeNodes: ITreeNode[] = profiles.map((e, i) => {
    return {
      id: i,
      hasCaret: false,
      icon: "document",
      label: e.label,
      isExpanded: false,
      isSelected: i === selected,
    };
  });

  return (
    <div className="row">
      <div className="col-xs-offset-1 col-xs-2">
        <H3>Profiles</H3>
        <Tree contents={treeNodes} onNodeClick={handleNodeClick} />
        <Divider style={{ marginTop: "20px", marginBottom: "20px" }} />
        <ControlGroup vertical fill>
          <InputGroup
            placeholder="New profile..."
            value={newProfile}
            onChange={(e) => {
              setNewProfile(e.target.value);
            }}
            rightElement={
              <Button minimal icon="add" onClick={handleNewProfile} />
            }
          />
          <Button
            icon="floppy-disk"
            disabled={isSaved}
            intent="primary"
            onClick={() => dispatch(save())}
          >
            Save
          </Button>
        </ControlGroup>
      </div>
      <div className="col-xs-8">
        {selected > -1 ? (
          <ProfileEditor index={selected} />
        ) : (
          <div>No profile selected. Select a profile or add a new one</div>
        )}
      </div>
    </div>
  );
}

export default Profile;
