#include "secret_input.h"

#include <FL/Fl_Secret_Input.H>

#include "event_handler.h"

class GSecret_Input : public EventHandler<Fl_Secret_Input> {
public:
  GSecret_Input(int x, int y, int w, int h, const char *label)
    : EventHandler<Fl_Secret_Input>(x, y, w, h, label) {}
};

GSecret_Input *go_fltk_new_Secret_Input(int x, int y, int w, int h, const char *text) {
  return new GSecret_Input(x, y, w, h, text);
}