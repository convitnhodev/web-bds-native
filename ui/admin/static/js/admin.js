document.addEventListener("DOMContentLoaded", function (event) {
    tinymce.init({
        selector: ".wysiwyg",
        // menubar: false,
        height: 1000,
        // statusbar: false,
        relative_urls: false,
        toolbar: 'media image link|bold italic | alignleft aligncenter alignright alignjustify | outdent indent',
        paste_block_drop: false,
        paste_merge_formats: true,
        smart_paste: true,
        remove_trailing_brs: true,
        valid_elements: "br,a[href|title|target|rel=nofollow noopener noreferrer],del,b,strong,del,i,blockquote,p[style],em,ul,li,ol,span[style],h1,h2,h3,h4,img[src|alt|width|style|class=width-fit|referrerpolicy=no-referrer|loading=lazy]",
        extended_valid_elements : "iframe[src|frameborder|style|scrolling|class|width|height|name|align]",
        valid_styles : { '*' : 'font-weight,font-style,text-decoration,text-align' },
        paste_preprocess: function(editor, args) {
            args.content = args.content.replaceAll("<br />", "")
            args.content = args.content.replaceAll("<br/>", "")
            args.content = args.content.replaceAll("<br>", "")
            args.content= args.content.replaceAll("<p></p>", "")
            args.content= args.content.replaceAll("<p>&nbsp;</p>", "")
        },
        plugins: 'media link image',
        link_rel_list: [
            { title: 'External Link', value: 'nofollow noopener noreferrer' },
            { title: 'OnPage Link', value: 'noopener' }
        ],
        link_target_list: [
            { title: 'New page', value: '_blank' },
            { title: 'Same page', value: '_self' }
          ]
    });
});