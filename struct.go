package zgen

import (
  "time"
)

// various generic constants
const (
  TimeFormatISOTZ   = time.RFC3339                // ISO format with Z timezone
  TimeFormatISOSTZ  = "2006-01-02 15:04:05 07:00" // ISO format with spaced timezone
  TimeFormatISO     = time.DateTime               // ISO time format with no timezone
  TimeFormatISODate = time.DateOnly               // ISO date format

  // DeepMerge flags
  FlagDeepMergeOverwriteDisabled = 0  // DeepMerge overwrite disabled
  FlagDeepMergePriorityFirst     = 1  // Set DeepMerge Priority on first element
  FlagDeepMergePrioritySecond    = 2  // Set DeepMerge Priority on second element
  FlagDeepMergeOverwriteEnabled  = 10 // DeepMerge overwrite enabled

  DeepCopyMaxIterationCount = 1000 // deepcopy will exit after this iteration count is reached to prevent infinite circular reference loops

  // ParserConfig Flags
  ParserModeNameOnly    = 1 // resulting map will only have struct property names as fields
  ParserModeTagsOnly    = 2 // resulting map will only have struct property tags as fields
  ParserModeNameAndTags = 3 // resulting map will have struct property names and tags as fields
  ParserModeNameIfNoTag = 4 // resulting map will have struct tags and names for the properties that have no tag as fields

)

type ParserConfig struct {
  EvaluateMethods bool     `json:"evaluate_methods"` // if true, it will try to determine if the struct is a Valuer or a Scanner for example and return its value instead of diving further
  KeepPointers    bool     `json:"keep_pointers"`    // keeps pointer values intact if true, dereferentiates them otherwise
  Mode            int      `json:"mode"`             // parses names and tags based on config value
  OmitEmpty       bool     `json:"omit_empty"`       // remove empty fields
  Tags            []string `json:"tags"`             // tag list to parse
}
