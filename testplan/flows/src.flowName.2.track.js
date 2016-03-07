function (input){
  for (i = 0; i < 4; i++){
    var searchParam = Math.floor((Math.random() * 100000));
    var lookup = Flow.Drop.findFHash("/loadtest/query/cities","" + searchParam);
    var key = "date" + i;
    input.elems[key] = lookup.creationDate;
  }
  var r = Math.floor((Math.random() * 5) + 1);
  input.elems.chosen = r;

  var p = "/loadtest/query/output/" + r;
  var output = {};
  output[p] = [input]
  return output;
}
