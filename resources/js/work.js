/*
  ~ The MIT License (MIT)
  ~
  ~ Copyright (c) 2016.  Georg Beier. All rights reserved.
  ~
  ~ Permission is hereby granted, free of charge, to any person obtaining a copy
  ~ of this software and associated documentation files (the "Software"), to deal
  ~ in the Software without restriction, including without limitation the rights
  ~ to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
  ~ copies of the Software, and to permit persons to whom the Software is
  ~ furnished to do so, subject to the following conditions:
  ~
  ~ The above copyright notice and this permission notice shall be included in all
  ~ copies or substantial portions of the Software.
  ~
  ~ THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
  ~ IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
  ~ FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
  ~ AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
  ~ LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
  ~ OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
  ~ SOFTWARE.
*/

    var actTab = 'enrol';
    var prevTab = '';
    var searchFlag = '';
    var bounceProtect = '';
    var searchStates = {};
    var finderValues = {};
    $(document).ready(function(){
        // bind function to #find button:
        $("#find").click(defaultFind);
        // bind function to tab selection event
        $('a[data-toggle="tab"]').on('shown.bs.tab', tabSwitchHandler);
        $('#do-cancellation').click(true, onCancellationRadioClick);
        $('#undo-cancellation').click(false, onCancellationRadioClick);
    });

    // default finder function for find button clicks
    // send query fields to server and load response into div #qapps
    function defaultFind() {
        var search1 = $("input#search1").val()
        var search2 = $("input#search2").val()
        var csrfid = $("input#csrf_id_find").val()
        var findUrl;
        if(typeof(finderValues[actTab]) !== 'undefined') {
            findUrl = finderValues[actTab].findUrl;
            $("#qapps").load(findUrl, {lastname: search1, firstname: search2, csrf_token: csrfid,
                action: actTab, flag: actTab === 'cancellation' ? searchFlag : ''},
                showSelect);
        }
    }

    // handler for switch tab events
    function tabSwitchHandler(e) {
        actTab = (""+e.target).replace(/.*#/, "");
        prevTab = (""+e.relatedTarget).replace(/.*#/, "");
        if(bounceProtect != prevTab + actTab) {
            bounceProtect = prevTab + actTab;
            if('' != prevTab) {
                searchStates[prevTab].saveContent();
            }
            searchStates[actTab].restoreContent();
        }
    }

    // handler for switching exmatriculation radio buttons
    function onCancellationRadioClick(e) {
        var val = e.data;
        if(val) {
            searchFlag = '';
        } else {
            searchFlag = 'undo';
        }
        if(actTab == 'cancellation') {
            $('#find').click();
        }
    }

    // select the target area and service url for active tab and load from server
    function showSelect() {
        var selid = $('#sresult').val();
        var csrfid = $("input#csrf_id_find").val()
        if(typeof(finderValues[actTab]) !== 'undefined' && selid > 0) {
            var target = finderValues[actTab].target;
            var url = finderValues[actTab].targetUrl;
            $(target).load(url, {appid: selid, action: actTab, csrf_token: csrfid,
            flag: actTab == 'cancellation' ? searchFlag : ''});
        }
    }

    // constructor function for objects that save and recover search results html when tabs are changed
    function SearchState(fid, s1id, s2id) {
        var selectid = fid;
        var search1 = s1id;
        var search2 = s2id;
        var content =  $(selectid).html();
        var sval1 = '';
        var sval2 = ''
        var idx = -1;
        this.saveContent = function() {
            content = $(selectid).html();
            if($(selectid +' select').length) {
                idx = $(selectid +' select').val();
            }
            sval1 = $(search1).val();
            sval2 = $(search2).val();
        };
        this.restoreContent = function() {
            $(selectid).html(content);
            /*
            if(search1.slice(1, 2) === '!') {
                search1 = search1.replace('!', '');
                console.log("replaced: " + search1);
            }
            if(search2.slice(1, 2) === '!') {
                search2 = search2.replace('!', '');
                console.log("replaced: " + search2);
            }
            */
            if($(selectid +' select').length) {
                $(selectid +' select').val(idx);
            }
            $(search1).val(sval1);
            $(search2).val(sval2);
        };
    }

    // post form data and check returned data. if it's not empty, edited object was
    // changed in the database. so normally update form fields and mark changes.
    // in case object was deleted, show alert and refresh find
    // if data is empty, call load area
    function postSubmissionAndReload(posturl, values, loadurl, area) {
        $.post(posturl, values, function( data ) {
            if( data ) {
                if (typeof data === 'string' || data instanceof String) {
                    // unexpected!
                    alert(data);
                    loadArea(loadurl, area);
                } else {
                    var keys = Object.keys(data);
                    keys.forEach(function(key) {
                        mergeVal(key, data[key], area);
                    });
                }
            } else {
                loadArea(loadurl, area);
            }
        });
    }

    // update fields with changed values and mark these fields
    function mergeVal( key, value, area ) {
        elements = Object.keys(value);
        // console.log("area: " + area + " key: " + key + " value: " + value['Other']);
        bg = ""
        setval = value['Other']
        /*
            NONE MergeDiffType = iota       // no changes at all
            MINE                            // only my value changed
            THEIRS                          // only their value changed
            BOTH                            // both values changed and are different
            SAME                            // both values changed but are equal
        */
        switch(value['Conflict']) {
            case 1:
                bg = "background-color: lightblue";
                setval = value['Mine']
                break;
            case 2:
                bg = "background-color: yellow";
                break;
            case 3:
                bg = "background-color: orangered";
                $(area + ' #' + key).prop('title', '[... ' + value['Mine'] + ' ?]')
                break;
            case 4:
                bg = "background-color: lightgreen";
                break;
        }
        $(area + ' #' + key + '-icon').attr('style', bg);
        $(area + ' #' + key).val(setval)
   }

    // load tab area with next from search field, removing previous from search results.
    // clear area if no further results in search results list.
    function loadArea(url, area) {
        var sid = $('#sresult').val();
        var idx = $('#sresult').prop('selectedIndex');
        var csrfid = $("#csrf_id_find").val()
        $('option[value="'+sid+'"]').remove();
        var l = $('#sresult').children('option').length;
        // goto next in list
        if(l > idx) {
            $('#sresult').prop('selectedIndex', idx);
        // else goto prev or empty
        } else {
            $('#sresult').prop('selectedIndex', idx - 1);
        }
        if(l > 0) {
            var selid = $('#sresult').val();
            $(area).load(url, {appid: selid, csrf_token: csrfid, action: actTab});
        } else {
            $(area).html("");
        }
    }

    // enable all input fields of current form, disable button
    // event.data should have to fields:
    // f: id of form with fields to be enabled
    // b: id of button to be disabled
    function enableFormFields(event) {
        aForm = event.data.f;
        aButton = event.data.b
        console.log('aButton: ' + aButton);
        $(aForm + ' input[disabled]').removeAttr("disabled");
        $(aForm + ' select[disabled]').removeAttr("disabled");
        if(aButton) {
            $(aButton).attr('disabled', 'disabled');
        }
    }
