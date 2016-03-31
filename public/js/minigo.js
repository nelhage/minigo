(function () {
   function longColor(c) {
     if(c == "W") {
       return "white";
     } else if(c == "B") {
       return "black";
     } else {
       return '';
     }
   }
   var size = 9;

   var GoSquare = React.createClass({
       render: function() {
           var classes = ["stone"];
           var c = longColor(this.props.contents);
           if (c == "") {
             classes.push("empty");
           } else {
             classes.push(c)
           }
           return (
             <div className="square">
               <div className={classes.join(" ")}></div>
               <div className="grid"></div>
             </div>
           );
       }
   });

   var GoBoard = React.createClass({
       getInitialState: function() {
         return {
           positions: {},
           to_move: 'W',
         };
       },
       componentDidMount: function() {
         $.ajax({
           url: "/board.json",
           dataType: 'json',
           cache: false,
           success: function(data) {
             this.setState(data);
           }.bind(this),
           error: function(xhr, status, err) {
             console.error("/board.json", status, err.toString());
           }.bind(this),
         });
       },
       at: function(r, c) {
         return this.state.positions[c+","+r];
       },
       render: function() {
         var rows = [];
         for (var r = 0; r < this.props.size; r++) {
           var row = [];
           for (var c = 0; c < this.props.size; c++) {
             row.push(<GoSquare key={c} contents={this.at(r,c)}/>);
           }
           rows.push(<div className="row" key={r}>{row}</div>);
         }
         var classes = ["goboard", longColor(this.state.to_move)];
         return (
           <div className={classes.join(" ")}>
             {rows}
           </div>
         );
       }
   });

   ReactDOM.render(
     <GoBoard size={size} />,
     document.getElementById('content')
   );
 })();
