var USERS_BY_ID = {};

const NOTE_TYPES = [
    'Marginalia',
    'Meta',
    'Prediction',
    'Question',
];

const classNamesByName = {
    noteTypeButton: 'note-type-button',
    primaryButton: 'mui-btn--primary',
};

const classesByName = {
    noteTypeButton: '.' + classNamesByName.noteTypeButton,
};

// CREATE ELEMENTS
const $createButtonWithText = function(text) {
    return $('<button>').addClass('mui-btn').text(text);
};

const activateButton = function($button) {
    $button.addClass(classNamesByName.primaryButton);
};

const deactivateButton = function($button) {
    $button.removeClass(classNamesByName.primaryButton);
};

const isButtonActive = function($button) {
    return $button.hasClass(classNamesByName.primaryButton);
};

const $createHalfRowDivWithElement = function($element) {
    return $('<div>').addClass('mui-col-xs-6').append($element);
};

const $createRowWithTwoElements = function($element1, $element2) {
    const $row = $('<div>').addClass('mui-row');

    const $column1 = $createHalfRowDivWithElement($element1);
    const $column2 = $createHalfRowDivWithElement($element2);

    return $row
        .append($column1)
        .append($column2);
};

const $createGridContainer = function() {
    return $('<div>').addClass('mui-container-fluid');
};

const $createTextarea = function(labelText) {
    const $textarea = $('<textarea>').prop('required', true).prop('rows', 4);
    const $label = $('<label>').text(labelText);

    return $('<div>').addClass('mui-textfield')
        .append($textarea)
        .append($label);
};

// NOTES
const $createAuthor = function(authorId) {
    const user = USERS_BY_ID[authorId];

    return $('<span>', {text: user.displayName});
};

const $createType = function(type) {
    return $('<span>', {text: type});
};

const $createCreationTime = function(creationTime) {
    return $('<span>', {text: moment(creationTime).fromNow()});
};

const $createContent = function(content) {
    return $('<div>', {text: content});
};

const $createDivider = function() {
    return $('<span>', {text: ' - '});
};

const $createNote = function(note) {
    const $author = $createAuthor(note.authorId);
    const $type = $createType(note.type);
    const $creationTime = $createCreationTime(note.creationTime);
    const $content = $createContent(note.content);

    const $header = $('<div>').addClass('note-header')
        .append($author).append($createDivider())
        .append($type).append($createDivider())
        .append($creationTime);

    return $('<div>').addClass('note')
        .append($header)
        .append($content);
};

// ADD NOTE
const $createAddNoteModal = function() {
    const $modal = $('<div>').addClass('modal');

    const $buttons = NOTE_TYPES.map(noteType => {
        return $createButtonWithText(noteType).addClass(classNamesByName.noteTypeButton);
    });

    const $textarea = $createTextarea('Note').addClass('note-content');

    return $modal
        .append($createGridContainer()
            .append($createRowWithTwoElements($buttons[0], $buttons[1]))
            .append($createRowWithTwoElements($buttons[2], $buttons[3])))
        .append($textarea);
};

const activateModal = function($modal) {
    mui.overlay('on', $modal.get(0));
};

$(function() {
    const $addNoteModal = $createAddNoteModal();

    $.get('/api/user', function(usersById) {
        USERS_BY_ID = usersById;

        $.get('/api/note', function(notes) {
            const $notes = $('#notes');

            notes.forEach((note) => {
                $notes.append(
                    $createNote(note)
                );
            });
        });
    });

    $('#add-note-button').click(function() {
        activateModal($addNoteModal);
    });

    $(document).on('click', classesByName.noteTypeButton, function() {
        const $clickedButton = $(this);

        if (isButtonActive($clickedButton)) {
            deactivateButton($clickedButton);
        } else {
            $(classesByName.noteTypeButton).each(function() {
                const $button = $(this);

                if ($button.is($clickedButton)) {
                    activateButton($button);
                } else {
                    deactivateButton($button);
                }
            });
        }
    });
});
