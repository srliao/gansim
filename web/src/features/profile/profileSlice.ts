import { createSlice, PayloadAction, createSelector } from "@reduxjs/toolkit";
import { Position, Toaster } from "@blueprintjs/core";
import {
  IProfile,
  STAT_TYPE_ATK,
  STAT_TYPE_ATKP,
  STAT_TYPE_CR,
  STAT_TYPE_ELEP,
  STAT_TYPE_HP,
} from "types";
import { defaultMods } from "./defaults";
import { AppThunk } from "app/store";
import { stat } from "fs";

interface IProfileState {
  profiles: IProfile[];
  saved: boolean;
  loaded: boolean;
}

const initialState: IProfileState = {
  profiles: [],
  saved: true,
  loaded: false,
};

//toaster class
const AppToaster = Toaster.create({
  className: "app-notification-toaster",
  position: Position.BOTTOM,
});

export function save(): AppThunk {
  return function (dispatch, getState) {
    dispatch(setSaved(true));
    var state = getState();
    console.log(state.profile);
    //clone the state
    localStorage.setItem("gansim-profile", JSON.stringify(state.profile));
    AppToaster.show({
      message: "Data saved successfully",
      intent: "success",
    });
  };
}

export function load(): AppThunk {
  return function (dispatch, getState) {
    const loadState = () => {
      try {
        const serializedState = localStorage.getItem("gansim-profile");
        if (serializedState === null) {
          return undefined;
        }
        return JSON.parse(serializedState);
      } catch (err) {
        return undefined;
      }
    };

    var state = loadState();
    console.log("reading from localStorage: ", state);
    if (state) {
      //load state
      console.log(state);
      dispatch(loadFromStorage(state));
    }
  };
}

export const profileSlice = createSlice({
  name: "profile",
  initialState,
  reducers: {
    setSaved: (state, action: PayloadAction<boolean>) => {
      state.saved = action.payload;
    },
    setLoaded: (state, action: PayloadAction<boolean>) => {
      state.loaded = action.payload;
    },
    loadFromStorage: (state, action: PayloadAction<IProfileState>) => {
      state.profiles = action.payload.profiles;
      state.saved = true;
      state.loaded = true;
    },
    pAdd: (state, action: PayloadAction<{ label: string }>) => {
      //check folder exist
      state.profiles.push({
        id: Date.now(),
        label: action.payload.label,
        character: {
          level: 0,
          base_atk: 0,
          mods: defaultMods(),
        },
        weapon: {
          base_atk: 0,
          mods: [],
        },
        enemy: {
          level: 100,
          phy_resist: 0.1,
          ele_resist: 0.1,
        },
        artifact_levels: 20,
        artifact_main_stats: {
          FLOWER: STAT_TYPE_HP,
          FEATHER: STAT_TYPE_ATK,
          SANDS: STAT_TYPE_ATKP,
          GOBLET: STAT_TYPE_ELEP,
          CIRCLET: STAT_TYPE_CR,
        },
        abilities: [],
      });
      state.saved = false;
    },
    pEdit: (state, action: PayloadAction<{ index: number; p: IProfile }>) => {
      if (
        action.payload.index > -1 &&
        action.payload.index < state.profiles.length
      ) {
        state.profiles[action.payload.index] = action.payload.p;
        state.saved = false;
      }
    },
  },
});

export const {
  pAdd,
  pEdit,
  setLoaded,
  setSaved,
  loadFromStorage,
} = profileSlice.actions;

export default profileSlice.reducer;
