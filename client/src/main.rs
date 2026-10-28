mod net;

use fltk::{app, prelude::*, window::Window};

fn main() {
    let app = app::App::default();
    let mut win = Window::new(100, 100, 1100, 700, "retrocode - M0");
    win.end();
    win.show();
    app.run().unwrap();
}
