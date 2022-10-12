import presetUno from "@unocss/preset-uno";
import presetWebFonts from "@unocss/preset-web-fonts";

export default {
  presets: [
    presetUno(),
    presetWebFonts({
      provider: "google",
      fonts: {
        sans: "Barlow",
      },
    }),
  ],
};
