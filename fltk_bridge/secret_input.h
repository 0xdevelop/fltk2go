#pragma once

#ifdef __cplusplus
extern "C" {
#endif

  typedef struct GSecret_Input GSecret_Input;

  extern GSecret_Input *go_fltk_new_Secret_Input(int x, int y, int w, int h, const char *text);

#ifdef __cplusplus
}
#endif