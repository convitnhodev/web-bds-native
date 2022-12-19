<<<<<<< HEAD
import presetTypography from '@unocss/preset-typography';
import presetUno from '@unocss/preset-uno';
import presetWebFonts from '@unocss/preset-web-fonts';
import svgToDataUri from 'mini-svg-data-uri';
=======
import presetUno from '@unocss/preset-uno';
import presetWebFonts from '@unocss/preset-web-fonts';
>>>>>>> 5ed448758ec77912d8f15b1cd516eb78d5d5ec71

export default {
  presets: [
    presetUno(),
    presetWebFonts({
      provider: 'google',
      fonts: {
        sans: 'Barlow',
      },
    }),
<<<<<<< HEAD
    presetTypography(),
=======
>>>>>>> 5ed448758ec77912d8f15b1cd516eb78d5d5ec71
  ],
  preflights: [
    {
      getCSS: () => '.preload *{transition-property: none;}',
    },
<<<<<<< HEAD
    {
      getCSS: ({ theme }) => `.dee-ol {
  padding-left: 4.75rem;
}
.dee-ol li {
  position: relative;
  min-height: 4rem;
}
.dee-ol li:where(:not(:first-child)) {
  margin-top: 2rem;
}
.dee-ol li::before {
  left: -4.75rem;
  position: absolute;
  font-weight: bold;
  content: counter(list-item);
  width: 4rem;
  height: 4rem;
  display: inline-flex;
  justify-content: center;
  align-items: center;
}
.dee-ol1 li::before {
  background-color: ${theme.colors.sky['600']};
  border-radius: ${theme.borderRadius.md};
  color: #fff;
  font-size: 2rem;
  justify-content: center;
  align-items: center;
}
.dee-ol2 li::before {
  color: ${theme.colors.sky['600']};;
  font-size: 3rem;
  justify-content: end;
  align-items: center;
}
`,
    },
    {
      getCSS: ({ theme }) => {
        const spacing = {
          px: '1px',
          0: '0px',
          0.5: '0.125rem',
          1: '0.25rem',
          1.5: '0.375rem',
          2: '0.5rem',
          2.5: '0.625rem',
          3: '0.75rem',
          3.5: '0.875rem',
          4: '1rem',
          5: '1.25rem',
          6: '1.5rem',
          7: '1.75rem',
          8: '2rem',
          9: '2.25rem',
          10: '2.5rem',
          11: '2.75rem',
          12: '3rem',
          14: '3.5rem',
          16: '4rem',
          20: '5rem',
          24: '6rem',
          28: '7rem',
          32: '8rem',
          36: '9rem',
          40: '10rem',
          44: '11rem',
          48: '12rem',
          52: '13rem',
          56: '14rem',
          60: '15rem',
          64: '16rem',
          72: '18rem',
          80: '20rem',
          96: '24rem',
        };
        const borderWidth = { DEFAULT: '1px' };

        const inputsClasses = [
          "[type='text']",
          "[type='email']",
          "[type='url']",
          "[type='password']",
          "[type='number']",
          "[type='date']",
          "[type='datetime-local']",
          "[type='month']",
          "[type='search']",
          "[type='tel']",
          "[type='time']",
          "[type='week']",
          '[multiple]',
          'textarea',
          'select',
        ];

        const rules = [
          {
            base: inputsClasses,
            class: [
              '.form-input',
              '.form-textarea',
              '.form-select',
              '.form-multiselect',
            ],
            styles: {
              appearance: 'none',
              'background-color': '#fff',
              'border-color': theme.colors.gray['500'],
              'border-width': borderWidth.DEFAULT,
              'border-radius': theme.borderRadius.none,
              'padding-top': spacing[2],
              'padding-right': spacing[3],
              'padding-bottom': spacing[2],
              'padding-left': spacing[3],
              'font-size': theme.fontSize.base[0],
              'line-height': theme.fontSize.base[0],
              '--un-shadow': '0 0 #0000',
            },
          },
          {
            base: inputsClasses.map((cssClass) => `${cssClass}:focus`),
            styles: {
              outline: '2px solid transparent',
              'outline-offset': '2px',
              '--un-ring-inset': 'var(--un-empty,/*!*/ /*!*/)',
              '--un-ring-offset-width': '0px',
              '--un-ring-offset-color': '#fff',
              '--un-ring-color': theme.colors.blue['600'],
              '--un-ring-offset-shadow':
                'var(--un-ring-inset) 0 0 0 var(--un-ring-offset-width) var(--un-ring-offset-color)',
              '--un-ring-shadow':
                'var(--un-ring-inset) 0 0 0 calc(1px + var(--un-ring-offset-width)) var(--un-ring-color)',
              'box-shadow':
                'var(--un-ring-offset-shadow), var(--un-ring-shadow), var(--un-shadow)',
              'border-color': theme.colors.blue['600'],
            },
          },
          {
            base: ['input::placeholder', 'textarea::placeholder'],
            class: ['.form-input::placeholder', '.form-textarea::placeholder'],
            styles: {
              color: theme.colors.gray['500'],
              opacity: '1',
            },
          },
          {
            base: ['::-webkit-datetime-edit-fields-wrapper'],
            class: ['.form-input::-webkit-datetime-edit-fields-wrapper'],
            styles: {
              padding: '0',
            },
          },
          {
            // Unfortunate hack until https://bugs.webkit.org/show_bug.cgi?id=198959 is fixed.
            // This sucks because users can't change line-height with a utility on date inputs now.
            // Reference: https://github.com/twbs/bootstrap/pull/31993
            base: ['::-webkit-date-and-time-value'],
            class: ['.form-input::-webkit-date-and-time-value'],
            styles: {
              'min-height': '1.5em',
            },
          },
          {
            // In Safari on macOS date time inputs are 4px taller than normal inputs
            // This is because there is extra padding on the datetime-edit and datetime-edit-{part}-field pseudo elements
            // See https://github.com/tailwindlabs/tailwindcss-forms/issues/95
            base: [
              '::-webkit-datetime-edit',
              '::-webkit-datetime-edit-year-field',
              '::-webkit-datetime-edit-month-field',
              '::-webkit-datetime-edit-day-field',
              '::-webkit-datetime-edit-hour-field',
              '::-webkit-datetime-edit-minute-field',
              '::-webkit-datetime-edit-second-field',
              '::-webkit-datetime-edit-millisecond-field',
              '::-webkit-datetime-edit-meridiem-field',
            ],
            class: [
              '.form-input::-webkit-datetime-edit',
              '.form-input::-webkit-datetime-edit-year-field',
              '.form-input::-webkit-datetime-edit-month-field',
              '.form-input::-webkit-datetime-edit-day-field',
              '.form-input::-webkit-datetime-edit-hour-field',
              '.form-input::-webkit-datetime-edit-minute-field',
              '.form-input::-webkit-datetime-edit-second-field',
              '.form-input::-webkit-datetime-edit-millisecond-field',
              '.form-input::-webkit-datetime-edit-meridiem-field',
            ],
            styles: {
              'padding-top': 0,
              'padding-bottom': 0,
            },
          },
          {
            base: ['select'],
            class: ['.form-select'],
            styles: {
              'background-image': `url("${svgToDataUri(
                `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20"><path stroke="${theme.colors.gray['500']}" stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 8l4 4 4-4"/></svg>`
              )}")`,
              'background-position': `right ${spacing[2]} center`,
              'background-repeat': `no-repeat`,
              'background-size': `1.5em 1.5em`,
              'padding-right': spacing[10],
              'print-color-adjust': `exact`,
            },
          },
          {
            base: ['[multiple]'],
            class: null,
            styles: {
              'background-image': 'initial',
              'background-position': 'initial',
              'background-repeat': 'unset',
              'background-size': 'initial',
              'padding-right': spacing[3],
              'print-color-adjust': 'unset',
            },
          },
          {
            base: [`[type='checkbox']`, `[type='radio']`],
            class: ['.form-checkbox', '.form-radio'],
            styles: {
              appearance: 'none',
              padding: '0',
              'print-color-adjust': 'exact',
              display: 'inline-block',
              'vertical-align': 'middle',
              'background-origin': 'border-box',
              'user-select': 'none',
              'flex-shrink': '0',
              height: spacing[4],
              width: spacing[4],
              color: theme.colors.blue['600'],
              'background-color': '#fff',
              'border-color': theme.colors.gray['500'],
              'border-width': borderWidth.DEFAULT,
              '--un-shadow': '0 0 #0000',
            },
          },
          {
            base: [`[type='checkbox']`],
            class: ['.form-checkbox'],
            styles: {
              'border-radius': theme.borderRadius.none,
            },
          },
          {
            base: [`[type='radio']`],
            class: ['.form-radio'],
            styles: {
              'border-radius': '100%',
            },
          },
          {
            base: [`[type='checkbox']:focus`, `[type='radio']:focus`],
            class: ['.form-checkbox:focus', '.form-radio:focus'],
            styles: {
              outline: '2px solid transparent',
              'outline-offset': '2px',
              '--un-ring-inset': 'var(--un-empty,/*!*/ /*!*/)',
              '--un-ring-offset-width': '2px',
              '--un-ring-offset-color': '#fff',
              '--un-ring-color': theme.colors.blue['600'],
              '--un-ring-offset-shadow': `var(--un-ring-inset) 0 0 0 var(--un-ring-offset-width) var(--un-ring-offset-color)`,
              '--un-ring-shadow': `var(--un-ring-inset) 0 0 0 calc(2px + var(--un-ring-offset-width)) var(--un-ring-color)`,
              'box-shadow': `var(--un-ring-offset-shadow), var(--un-ring-shadow), var(--un-shadow)`,
            },
          },
          {
            base: [`[type='checkbox']:checked`, `[type='radio']:checked`],
            class: ['.form-checkbox:checked', '.form-radio:checked'],
            styles: {
              'border-color': `transparent`,
              'background-color': `currentColor`,
              'background-size': `100% 100%`,
              'background-position': `center`,
              'background-repeat': `no-repeat`,
            },
          },
          {
            base: [`[type='checkbox']:checked`],
            class: ['.form-checkbox:checked'],
            styles: {
              'background-image': `url("${svgToDataUri(
                `<svg viewBox="0 0 16 16" fill="white" xmlns="http://www.w3.org/2000/svg"><path d="M12.207 4.793a1 1 0 010 1.414l-5 5a1 1 0 01-1.414 0l-2-2a1 1 0 011.414-1.414L6.5 9.086l4.293-4.293a1 1 0 011.414 0z"/></svg>`
              )}")`,
            },
          },
          {
            base: [`[type='radio']:checked`],
            class: ['.form-radio:checked'],
            styles: {
              'background-image': `url("${svgToDataUri(
                `<svg viewBox="0 0 16 16" fill="white" xmlns="http://www.w3.org/2000/svg"><circle cx="8" cy="8" r="3"/></svg>`
              )}")`,
            },
          },
          {
            base: [
              `[type='checkbox']:checked:hover`,
              `[type='checkbox']:checked:focus`,
              `[type='radio']:checked:hover`,
              `[type='radio']:checked:focus`,
            ],
            class: [
              '.form-checkbox:checked:hover',
              '.form-checkbox:checked:focus',
              '.form-radio:checked:hover',
              '.form-radio:checked:focus',
            ],
            styles: {
              'border-color': 'transparent',
              'background-color': 'currentColor',
            },
          },
          {
            base: [`[type='checkbox']:indeterminate`],
            class: ['.form-checkbox:indeterminate'],
            styles: {
              'background-image': `url("${svgToDataUri(
                `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 16 16"><path stroke="white" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h8"/></svg>`
              )}")`,
              'border-color': `transparent`,
              'background-color': `currentColor`,
              'background-size': `100% 100%`,
              'background-position': `center`,
              'background-repeat': `no-repeat`,
            },
          },
          {
            base: [
              `[type='checkbox']:indeterminate:hover`,
              `[type='checkbox']:indeterminate:focus`,
            ],
            class: [
              '.form-checkbox:indeterminate:hover',
              '.form-checkbox:indeterminate:focus',
            ],
            styles: {
              'border-color': 'transparent',
              'background-color': 'currentColor',
            },
          },
          {
            base: [`[type='file']`],
            class: null,
            styles: {
              background: 'unset',
              'border-color': 'inherit',
              'border-width': '0',
              'border-radius': '0',
              padding: '0',
              'font-size': 'unset',
              'line-height': 'inherit',
            },
          },
          {
            base: [`[type='file']:focus`],
            class: null,
            styles: {
              outline: `1px solid ButtonText , 1px auto -webkit-focus-ring-color`,
            },
          },
        ];

        const createStyleObject = ([key, value]) => {
          if (typeof value === 'object') {
            return Object.entries(value)
              .map((styles) => createStyleObject(styles))
              .join('\n');
          }

          return `${key}: ${value};`;
        };

        const style = rules.map((rule) => {
          const selector = rule.base.join(', ');

          const styles = Object.entries(rule.styles)
            .map((style) => createStyleObject(style))
            .join('\n');

          return `${selector} { ${styles} }`;
        });

        return style.join('\n');
      },
    },
=======
>>>>>>> 5ed448758ec77912d8f15b1cd516eb78d5d5ec71
  ],
};
